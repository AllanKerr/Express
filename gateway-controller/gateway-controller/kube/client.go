package kube

import (
	"k8s.io/client-go/kubernetes"
	"flag"
	"path/filepath"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// A wrapper around the Kubernetes go-client library client structure
// required to interact with the Kubernetes API.
type Client struct {
	*kubernetes.Clientset
}

/***************************************************************************************
*    Authenticating outside the cluster
*    Author: Marc Sluiter, David Xia, and Ahmet Alp Balkan
*    Date: Dec. 29, 2017
*    Code version: 6.0.0
*    Availability: https://github.com/kubernetes/client-go/blob/master/examples/out-of-cluster-client-configuration/main.go
*
***************************************************************************************/

// Creates a new client for interfacing with the Kubernetes API
// Kubernetes must be setup on the local system or the kubeconfig will not be found
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

// The location of the home directory to find the kubeconfig
func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}


// List the set of applications that have been deployed using the Express deploy command
func (client *Client) ListApplications(namespace string) ([]Application, error) {

	sClient := client.CoreV1().Services(namespace)
	services, err := sClient.List(metav1.ListOptions{
		LabelSelector: "group=services",
	})
	if err != nil {
		return nil, err
	}
	var applications []Application
	for _, service := range services.Items {
		applications = append(applications, &DefaultApplication{
			name: service.GetName(),
			port: service.Spec.Ports[0].Port,
			creation: service.GetCreationTimestamp().Time,
		})
	}
	return applications, nil
}