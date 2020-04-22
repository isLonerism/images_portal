package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	oc "github.com/bsuro10/images_portal/images-portal-grpc-client/interfaces/openshift_client"
	"github.com/bsuro10/images_portal/images-portal-grpc-server/api/docker"
	"google.golang.org/grpc"
)

type ProjectsRequest struct {
	APIEndpoint string
	Token       string
}

type ProjectsResponse struct {
	ProjectList []string
}

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

type PushResponse struct {
	message string
}

var (
	currentValue, currentSeconds = 0, -1
	mux                          sync.Mutex
)

func main() {
	http.HandleFunc("/load", load)
	http.HandleFunc("/push", push)
	http.HandleFunc("/projects", projects)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func projects(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var request ProjectsRequest
	if succeeded := getProjectsRequest(w, r, &request); !succeeded {
		log.Println("could not unmarshal request")
		return
	}

	projectList := getProjectsList(request)
	if projectList == nil {
		log.Println("could not generate project list")
		return
	}

	sendProjectsResponse(w, projectList)
}

func load(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var request LoadRequest
	if succeeded := getLoadRequest(w, r, &request); !succeeded {
		return
	}

	podInterface, err := deployPod(w)
	if err != nil {
		return
	}

	imageList := callGRPCLoad(w, podInterface, request)
	if imageList == nil {
		return
	}

	sendLoadResponse(w, podInterface, imageList)

}

func push(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var request PushRequest
	if succeeded := getPushRequest(w, r, &request); !succeeded {
		return
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

	conn, err := makeGRPCConnection(w, podInterface)
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

	message, err := docker.NewDockerClient(conn).TagAndPush(context.Background(), &a)
	if err != nil {
		log.Println(err)
		http.Error(w, "grpc server failed during push", http.StatusBadRequest)
	}

	log.Println(message)

	response := PushResponse{
		message: "successfully pushed",
	}
	sendJSONResponse(w, response)

	err = oc.DeletePod(podInterface.PodName)
	if err != nil {
		log.Println(err)
	}
	log.Println("finished handling pod " + podInterface.PodName)

}

func getProjectsList(request ProjectsRequest) []string {
	req, err := http.NewRequest("GET", request.APIEndpoint+"/apis/project.openshift.io/v1/projects", nil)
	if err != nil {
		http.Error(w, "could not create projects request", http.StatusInternalServerError)
		log.Println(err)
		return nil
	}

	req.Header.Add("Authorization", "Bearer "+request.Token)
	req.Header.Add("Accept", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		http.Error(w, "error response received from API", http.StatusBadRequest)
		log.Println(err)
		return nil
	}
	defer res.Body.Close()

	resJSON := map[string]interface{}{}
	err = json.NewDecoder(res.Body).Decode(&resJSON)
	if err != nil {
		http.Error(w, "unexpected response received from API", http.StatusInternalServerError)
		log.Println(err)
		return nil
	}

	var projectList []string
	for _, item := range resJSON["items"].([]interface{}) {
		projectList = append(projectList, item.(map[string]interface{})["metadata"].(map[string]interface{})["name"].(string))
	}

	return projectList
}

func callGRPCLoad(w http.ResponseWriter, podInterface oc.PodInterface, request LoadRequest) *docker.ImagesList {
	conn, err := makeGRPCConnection(w, podInterface)
	if err != nil {
		return nil
	}
	defer conn.Close()

	s3Object := docker.S3Object{
		S3Key:       request.S3Key,            //from client
		S3Bucket:    os.Getenv("S3Bucket"),    //env var
		S3Accesskey: os.Getenv("S3Accesskey"), //env var
		S3Secretkey: os.Getenv("S3Secretkey"), //env var
		S3Endpoint:  os.Getenv("S3Endpoint"),  //env var
		S3Region:    os.Getenv("S3Region"),    //env var
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

func makeGRPCConnection(w http.ResponseWriter, podInterface oc.PodInterface) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(podInterface.PodIP+":"+os.Getenv("GRPC_PORT"), grpc.WithInsecure()) //ip is from podInterface, port is env var
	if err != nil {
		log.Println(err)
		http.Error(w, "unable to open grpc connection", http.StatusInternalServerError)
	}
	return conn, err
}

func sendProjectsResponse(w http.ResponseWriter, projectList []string) {
	projectsResponse := ProjectsResponse{
		ProjectList: projectList,
	}

	sendJSONResponse(w, projectsResponse)
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

func getProjectsRequest(w http.ResponseWriter, r *http.Request, request *ProjectsRequest) bool {
	err := getRequest(w, r, request, "projects")
	if err != nil {
		return false
	}
	log.Println(request)

	if valid := isValidProjectsRequest(*request); !valid {
		http.Error(w, "missing projects request params", http.StatusBadRequest)
		return false
	}
	return true
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

func isValidProjectsRequest(request ProjectsRequest) bool {
	return request.APIEndpoint != "" && request.Token != ""
}

func isValidLoadRequest(request LoadRequest) bool {
	return request.S3Key != ""
}

func isValidPushRequest(request PushRequest) bool {
	for _, item := range request.Images.Images {
		if item.OldImage.Name == "" || item.NewImage.Name == "" {
			return false
		}
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

func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}
