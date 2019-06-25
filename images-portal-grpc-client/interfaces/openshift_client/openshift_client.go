package openshift_client

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gopkg.in/yaml.v2"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeyaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type PodInterface struct {
	PodName string
	PodIP   string
}

var (
	config    *rest.Config
	clientset *kubernetes.Clientset
)

func init() {
	var err error

	// Out-Of-Cluster configuration
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// In-Cluster configuration
	// config, err = rest.InClusterConfig()
	// if err != nil {
	// 	panic(err.Error())
	// }

	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
}

func DeployPod(value int) (PodInterface, error) {
	log.Println("value: " + strconv.Itoa(value))
	template, err := os.Open("/home/paas/pod.yml") //file is configmap, path is env var
	if err != nil {
		log.Fatal(err)
	}
	defer template.Close()

	templateData, err := ioutil.ReadAll(template)

	m := make(map[interface{}]interface{})
	yaml.Unmarshal([]byte(templateData), &m)

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	metadata := m["metadata"].(map[interface{}]interface{})
	labels := metadata["labels"].(map[interface{}]interface{})

	name := time.Now().Format("20060102-150405") + "-" + strconv.Itoa(value)

	metadata["name"] = name
	labels["name"] = name

	newYaml, err := yaml.Marshal(&m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	tempPath := filepath.Join(os.TempDir(), name)
	podYaml, err := os.Create(tempPath)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer os.Remove(tempPath)
	_, err = podYaml.Write(newYaml)

	var pod apiv1.Pod

	podYaml.Seek(0, 0)
	kubeyaml.NewYAMLOrJSONDecoder(podYaml, 1024).Decode(&pod) //buffer size is env var

	fmt.Println("Creating pod...")
	result, err := clientset.CoreV1().Pods("myproject").Create(&pod) //project name is env var
	if err != nil {
		log.Println(err)
		return PodInterface{"", ""}, err
	}
	fmt.Printf("Created pod %q.\n", result.GetObjectMeta().GetName())

	watcher, err := clientset.CoreV1().Pods("myproject").Watch(metav1.ListOptions{ //project name is env var
		LabelSelector: "name=" + name,
	})
	if err != nil {
		log.Println(err)
		return PodInterface{"", ""}, err
	}

	ch := watcher.ResultChan()

	timer := time.NewTimer(20 * time.Second) //duration is env var
	go func() {
		<-timer.C
		watcher.Stop()
	}()
	log.Println(name)
	for pods := range ch {
		pod, ok := pods.Object.(*v1.Pod)
		if !ok {
			log.Println("unexpected type")
		}

		switch pod.Status.Phase {
		case "Failed":
			log.Println("Pod Failed")
			return PodInterface{"", ""}, errors.New("Pod failed")
		case "Running":
			log.Println("Pod Running")
			return PodInterface{
				PodName: name,
				PodIP:   pod.Status.PodIP,
			}, nil
		}
	}

	return PodInterface{"", ""}, errors.New("Pod " + name + " is Pending Forever")
}

func DeletePod(name string) error {
	return clientset.CoreV1().Pods("myproject").Delete(name, &metav1.DeleteOptions{}) //project name is env var
}

func homeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Println("no home dir")
		return ""
	}
	return home
}
