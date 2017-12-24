package kube

import (
	"k8s.io/client-go/kubernetes"
	"flag"
	"path/filepath"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/api/core/v1"
)

type Client struct {
	*kubernetes.Clientset
}

func NewDefaultClient() (*Client, error) {

	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &Client{clientset}, nil
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func (client *Client) ListServices(namespace string) ([]apiv1.Service, error) {

	sClient := client.CoreV1().Services(namespace)
	services, err := sClient.List(metav1.ListOptions{
		LabelSelector: "group=services",
	})
	return services.Items, err
}