package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	oc "github.com/bsuro10/images-portal/images-portal-grpc-client/interfaces/openshift_client"
	"github.com/bsuro10/images-portal/images-portal-grpc-server/api/docker"
	"google.golang.org/grpc"
)

type LoadRequest struct {
	S3Key string
}

type LoadResponse struct {
	Token  string
	Images []*docker.Image
}

type PushRequest struct {
	PodToken    string
	DockerToken string
	Images      docker.TagImagesList
}

var (
	currentValue, currentSeconds = 0, -1
	mux                          sync.Mutex
)

func main() {
	http.HandleFunc("/upload", uploadFileHandler)

	http.HandleFunc("/load", load)
	http.HandleFunc("/push", push)

	http.HandleFunc("/test", test)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func load(w http.ResponseWriter, r *http.Request) {
	var request LoadRequest
	if succeeded := getLoadRequest(w, r, &request); !succeeded {
		return
	}

	podInterface, err := deployPod(w)
	if err != nil {
		return
	}

	timer := time.NewTimer(20 * time.Second)

	<-timer.C

	log.Println("starting grpc")

	imageList := callGRPCLoad(w, request)
	if imageList == nil {
		return
	}

	sendLoadResponse(w, podInterface, imageList)

}

func push(w http.ResponseWriter, r *http.Request) {
	var request PushRequest
	if succeeded := getPushRequest(w, r, &request); !succeeded {
		return
	}

	conn, err := makeGRPCConnection(w)
	if err != nil {
		return
	}
	defer conn.Close()

	a := docker.TagAndPushObject{
		TagImages: &request.Images,
		AuthConfig: &docker.AuthConfig{ //get username and password from client
			Password: request.DockerToken,
			Username: "unused",
		},
	}

	podBytes, err := base64.StdEncoding.DecodeString(request.PodToken)
	if err != nil {
		log.Println(err)
		http.Error(w, "the session token is invalid", http.StatusBadRequest)
		return
	}
	var podInterface oc.PodInterface
	err = json.Unmarshal(podBytes, &podInterface)
	if err != nil {
		log.Println(err)
		http.Error(w, "the session token is invalid", http.StatusBadRequest)
		return
	}

	message, err := docker.NewDockerClient(conn).TagAndPush(context.Background(), &a)
	if err != nil {
		log.Println(err)
		http.Error(w, "grpc server failed during push", http.StatusBadRequest)
	}

	log.Println(message)

	// send message or error to client

	timer := time.NewTimer(60 * time.Second)

	<-timer.C
	err = oc.DeletePod(podInterface.PodName)
	if err != nil {
		log.Println(err)
	}
	log.Println("finished handling pod " + podInterface.PodName)

}

func callGRPCLoad(w http.ResponseWriter, request LoadRequest) *docker.ImagesList {
	conn, err := makeGRPCConnection(w)
	if err != nil {
		return nil
	}
	defer conn.Close()

	s3Object := docker.S3Object{
		S3Key:       request.S3Key,                              //from client
		S3Bucket:    "test",                                     //env var
		S3Accesskey: "2SXSSELJRNSSX9V4UYV8",                     //env var
		S3Secretkey: "p4wyUtqxYzj+1CeTJ8euxiwURJMr6swHPHEhj1gF", //env var
		S3Endpoint:  "http://10.0.0.4:9000",                     //env var
		S3Region:    "us-east-1",                                //env var
	}
	imageList, err := docker.NewDockerClient(conn).Load(context.Background(), &s3Object)
	if err != nil {
		log.Println(err)
		http.Error(w, "grpc server failed during load", http.StatusInternalServerError)
		return nil
	}

	log.Println(imageList)

	return imageList
}

func makeGRPCConnection(w http.ResponseWriter) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial("192.168.42.100:31192", grpc.WithInsecure()) //ip is from podInterface, port is env var
	if err != nil {
		log.Println(err)
		http.Error(w, "unable to open grpc connection", http.StatusInternalServerError)
	}
	return conn, err
}

func sendLoadResponse(w http.ResponseWriter, podInterface oc.PodInterface, imageList *docker.ImagesList) {
	podBytes, err := json.Marshal(podInterface)
	if err != nil {
		log.Println(err)
		http.Error(w, "cannot create token", http.StatusInternalServerError)
	}

	loadResponse := LoadResponse{
		Token:  base64.StdEncoding.EncodeToString(podBytes),
		Images: imageList.Images,
	}

	log.Println(loadResponse)

	sendJSONResponse(w, loadResponse)
}

func deployPod(w http.ResponseWriter) (oc.PodInterface, error) {
	podInterface, err := oc.DeployPod(getValue())
	if err != nil {
		log.Println(err)
		return oc.PodInterface{}, err
	}
	log.Println(podInterface)
	return podInterface, nil
}

func getLoadRequest(w http.ResponseWriter, r *http.Request, request *LoadRequest) bool {
	err := getRequest(w, r, request, "load")
	if err != nil {
		return false
	}
	log.Println(request)

	if valid := isValidLoadRequest(*request); !valid {
		http.Error(w, "missing load request params", http.StatusBadRequest)
		return false
	}
	return true
}

func getPushRequest(w http.ResponseWriter, r *http.Request, request *PushRequest) bool {
	err := getRequest(w, r, request, "push")
	if err != nil {
		return false
	}
	log.Println(request)

	if valid := isValidPushRequest(*request); !valid {
		http.Error(w, "missing push request params", http.StatusBadRequest)
		return false
	}
	return true
}

func isValidLoadRequest(request LoadRequest) bool {
	return request.S3Key != ""
}

func isValidPushRequest(request PushRequest) bool {
	log.Println(request.PodToken)
	log.Println(request.DockerToken)
	log.Println(request.Images)
	for _, item := range request.Images.Images {
		log.Println(item.OldImage)
		log.Println(item.NewImage)
		if item.OldImage.Name == "" || item.NewImage.Name == "" {
			return false
		}
	}
	var r docker.TagImagesList
	for _, item := range r.Images {
		log.Println(item.OldImage.Name)
	}
	return request.PodToken != "" && request.DockerToken != ""
}

func getRequest(w http.ResponseWriter, r *http.Request, request interface{}, urlPath string) error {
	reqBody, err := getBody(w, r, urlPath)
	if err != nil {
		return err
	}

	return getJSONRequest(w, reqBody, request, urlPath)
}

func getBody(w http.ResponseWriter, r *http.Request, urlPath string) ([]byte, error) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "cannot read "+urlPath+"request body", http.StatusBadRequest)
		return nil, err
	}
	return reqBody, nil
}

func getJSONRequest(w http.ResponseWriter, reqBody []byte, request interface{}, urlPath string) error {
	err := json.Unmarshal(reqBody, &request)
	if err != nil {
		log.Println(err)
		http.Error(w, "cannot unmarshal "+urlPath+" request body", http.StatusBadRequest)
		return err
	}
	return nil
}

func sendJSONResponse(w http.ResponseWriter, message interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(message)
	if err != nil {
		log.Println(message)
		log.Println(err)
		http.Error(w, "cannot encode response to json", http.StatusInternalServerError)
	}
}

func getValue() int {
	mux.Lock()
	value := calcValue()
	mux.Unlock()
	return value
}

func calcValue() int {
	_, _, seconds := time.Now().Clock()
	if seconds != currentSeconds {
		currentSeconds = seconds
		currentValue = 0
	} else {
		currentValue++
	}
	return currentValue
}

func test(w http.ResponseWriter, r *http.Request) {

	var request PushRequest
	if succeeded := getPushRequest(w, r, &request); !succeeded {
		return
	}

	log.Println("after getPushRequest")

	var podInterface oc.PodInterface
	podBytes, err := base64.StdEncoding.DecodeString(request.PodToken)
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal(podBytes, &podInterface)
	if err != nil {
		log.Println(err)
	}

	log.Println("finished")

	log.Println(podInterface)
}

func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	mux.Lock()
	_, _, seconds := time.Now().Clock()
	if seconds != currentSeconds {
		currentSeconds = seconds
		currentValue = 0
	} else {
		currentValue++
	}
	value := currentValue
	mux.Unlock()
	log.Println("seconds: " + strconv.Itoa(seconds))
	handleUpload(value, w, r)
}

func handleUpload(value int, w http.ResponseWriter, r *http.Request) {

	podInterface, err := oc.DeployPod(value)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(podInterface.PodName + " : " + podInterface.PodIP)

	h := sha256.New()

	h.Write([]byte(podInterface.PodName + podInterface.PodIP))
	bh := base64.StdEncoding.EncodeToString(h.Sum(nil))

	log.Println(bh)

	retPodInterface := podInterface

	conn, err := grpc.Dial("192.168.42.100:31190", grpc.WithInsecure()) //ip is from podInterface, port and security are env vars
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()
	c := docker.NewDockerClient(conn)

	ctx := context.Background()

	s3Object := docker.S3Object{
		S3Key:       "statsd-exporter-vault.tar",                //from client
		S3Bucket:    "test",                                     //env var
		S3Accesskey: "2SXSSELJRNSSX9V4UYV8",                     //env var
		S3Secretkey: "p4wyUtqxYzj+1CeTJ8euxiwURJMr6swHPHEhj1gF", //env var
		S3Endpoint:  "http://10.0.0.4:9000",                     //env var
		S3Region:    "us-east-1",                                //env var
	}
	images, err := c.Load(ctx, &s3Object)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(images)

	// loadResponse := LoadResponse{
	// 	Token:  bh,
	// 	Images: images,
	// }

	// log.Println(loadResponse)

	// err = sendJSONResponse(w, loadResponse)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	//send bh and images to client

	tag1 := docker.TagImage{
		OldImage: &docker.Image{
			Name: "prom/statsd-exporter:v0.5.0",
		},
		NewImage: &docker.Image{
			Name: "docker-registry.default.svc:5000/myproject/prom:prom-tag",
		},
	}

	tag2 := docker.TagImage{
		OldImage: &docker.Image{
			Name: "quay.io/coreos/vault:0.9.1-0",
		},
		NewImage: &docker.Image{
			Name: "docker-registry.default.svc:5000/myproject/vault:vault-tag",
		},
	}

	tags := docker.TagImagesList{
		Images: []*docker.TagImage{
			&tag1,
			&tag2,
		},
	}

	a := docker.TagAndPushObject{
		TagImages: &tags,
		AuthConfig: &docker.AuthConfig{ //get username and password from client
			Password: "eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJteXByb2plY3QiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlY3JldC5uYW1lIjoiYnVpbGRlci10b2tlbi1uZHFxYyIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJidWlsZGVyIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQudWlkIjoiYWUzNThmMDUtNTYwOS0xMWU5LTgzMTMtNTI1NDAwMWEwNmQxIiwic3ViIjoic3lzdGVtOnNlcnZpY2VhY2NvdW50Om15cHJvamVjdDpidWlsZGVyIn0.z5lZSat6hA2h-bqSyAuWHD_xLUYE_itXQQw48xoruMxHb3JKTTZfKQwaDByQpPbaCBubqvtcFt3xRLEE7HwB7lF-xpODa-K5U0glpLALFuIg8_MNshO9mSvMGrxqzXdW3ZtfvB5SjrtqMX5gKg3V8zRP2dwci4_snPovHaMduq7gHX0p-fNmER3rKInOVx7KCV1VfntRJuOTSskmely_wMDq4jL6JK0j_Py-CLHe1w9Cvnu87UUaahj5gRtV_OLxY_9w4azhLuInaj4lBdeQqDKMHq7LDkdYICt1zrdQKAzTM5qfivEVwcZH8C-6ui2GBXu626HuMQh6e2ANV9T3Cw",
			Username: "unused",
		},
	}

	message, err := c.TagAndPush(ctx, &a)
	if err != nil {
		log.Println(err)
	}

	log.Println(message)

	// send message or error to client

	err = oc.DeletePod(retPodInterface.PodName)
	if err != nil {
		log.Println(err)
	}
	log.Println("finished handling pod " + retPodInterface.PodName)
	//}()
}

func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}
