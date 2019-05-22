package openshift_client

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

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

func DeployPod() (string, error) {
	pod := &apiv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "pod-" + time.Now().Format("20060102-150405"),
			Labels: map[string]string{
				"name": "pod-" + time.Now().Format("20060102-150405"),
			},
		},
		Spec: apiv1.PodSpec{
			Containers: []apiv1.Container{
				{
					Name:  "grpc-server",
					Image: "bsuro10/images-portal-grpc-server:0.1",
					Ports: []apiv1.ContainerPort{
						{
							Name:          "http",
							Protocol:      apiv1.ProtocolTCP,
							ContainerPort: 7777,
						},
					},
				},
			},
			RestartPolicy: "Never",
		},
	}

	fmt.Println("Creating pod...")
	result, err := clientset.CoreV1().Pods("myproject").Create(pod)
	if err != nil {
		log.Println(err)
		return "", err
	}
	fmt.Printf("Created pod %q.\n", result.GetObjectMeta().GetName())

	watcher, err := clientset.CoreV1().Pods("myproject").Watch(metav1.ListOptions{
		LabelSelector: "name=" + result.GetObjectMeta().GetName(),
	})
	if err != nil {
		log.Println(err)
		return "", err
	}

	ch := watcher.ResultChan()

	timer := time.NewTimer(20 * time.Second)
	go func() {
		<-timer.C
		watcher.Stop()
	}()
	for pods := range ch {
		pod, ok := pods.Object.(*v1.Pod)
		if !ok {
			log.Println("unexpected type")
		}

		switch pod.Status.Phase {
		case "Failed":
			log.Println("Pod Failed")
			return "", errors.New("Pod failed")
		case "Running":
			log.Println("Pod Running")
			return pod.Status.PodIP, nil
		}
	}

	log.Println("Finished watch")
	return "", errors.New("Pod Pending Forever")
}

func DeletePod() {
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
