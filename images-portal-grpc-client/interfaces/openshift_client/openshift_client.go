package openshift_client

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	apiv1 "k8s.io/api/core/v1"
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

func DeployPod() {
	pod := &apiv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "web",
		},
		Spec: apiv1.PodSpec{
			Containers: []apiv1.Container{
				{
					Name:  "web",
					Image: "nginx:1.12",
					Ports: []apiv1.ContainerPort{
						{
							Name:          "http",
							Protocol:      apiv1.ProtocolTCP,
							ContainerPort: 80,
						},
					},
				},
			},
		},
	}

	fmt.Println("Creating pod...")
	result, err := clientset.CoreV1().Pods("myproject").Create(pod)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("Created pod %q.\n", result.GetObjectMeta().GetName())
}

func DeletePod() {
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
