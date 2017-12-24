package handlers

import (
	"github.com/spf13/cobra"
	"fmt"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"gateway-controller/kube"
	"k8s.io/apimachinery/pkg/api/errors"
	"io/ioutil"
)

type DeployHandler struct {
	name string
	client *kube.Client
	transactions []kube.Transaction
	err error
}

func NewDeployHandler(client *kube.Client, name string) *DeployHandler {
	return &DeployHandler{
		name: name,
		client: client,
	}
}

func (handler *DeployHandler) Execute(txn kube.Transaction, obj interface{}) error {

	if handler.err != nil {
		return handler.err
	}
	handler.err = txn.Execute(obj)
	if handler.err != nil {
		return handler.err
	}
	handler.transactions = append(handler.transactions, txn)
	return nil
}

func (handler *DeployHandler) Rollback() {

	for _, txn := range handler.transactions {
		if err := txn.Rollback(handler.name); err != nil {
			fmt.Printf("failed rollback: %v\n", err)
		}
	}
}

func (handler *DeployHandler) createService(name string, port int32) {

	service := kube.DefaultServiceConfig()

	service.ObjectMeta.Name = name
	service.ObjectMeta.Labels = map[string]string{
		"app": name,
		"group": "services",
	}

	service.Spec.Ports[0].Port = port
	service.Spec.Ports[0].TargetPort = intstr.FromInt(int(port))

	service.Spec.Selector = map[string]string {
		"app": name,
	}

	txn := kube.NewServiceTransaction(handler.client, apiv1.NamespaceDefault)
	if err := handler.Execute(txn, service); err == nil {
		fmt.Printf("Created service %q.\n", name)
	}
}

func (handler *DeployHandler) createDeployment(name string, image string, port int32, n int32) {

	labels := map[string]string{
		"app": name,
	}

	deployment := &appsv1beta2.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: labels,
		},
		Spec: appsv1beta2.DeploymentSpec{
			Replicas: &n,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: name,
					Labels: labels,
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  name,
							Image: image,
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: port,
								},
							},
							Lifecycle: &apiv1.Lifecycle{
								PreStop: &apiv1.Handler{
									Exec: &apiv1.ExecAction{
										Command: []string{"sleep", "15"},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	txn := kube.NewDeploymentTransaction(handler.client, apiv1.NamespaceDefault)
	if err := handler.Execute(txn, deployment); err == nil {
		fmt.Printf("Created deployment %q.\n", name)
	}
}

func (handler *DeployHandler) createAutoscaler(name string, min int32, max int32) {

	autoscaler := kube.DefaultAutoscalerConfig()

	autoscaler.Name = name
	autoscaler.Labels = map[string]string{
		"app": name,
	}

	autoscaler.Spec.ScaleTargetRef.Name = name
	autoscaler.Spec.MinReplicas = &min
	autoscaler.Spec.MaxReplicas = max

	txn := kube.NewAutoscalerTransaction(handler.client, apiv1.NamespaceDefault)
	if err := handler.Execute(txn, autoscaler); err == nil {
		fmt.Printf("Created autoscaler %q.\n", name)
	}
}

func (handler *DeployHandler) createEndpoints(name string, port int32, configFile string) {

	// If no config file is given, treat it as a private service with no endpoints
	if configFile == "" {
		return
	}
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		handler.err = err
		return
	}

	ingresses, err := kube.ParseConfig(name, port, file)
	if err != nil {
		handler.err = err
		return
	}

	hasError := handler.err != nil

	txn := kube.NewIngressTransaction(handler.client, apiv1.NamespaceDefault)
	if err := handler.Execute(txn, ingresses); err == nil {
		fmt.Printf("Created ingresses %q.\n", name)
	} else if !hasError {
		// Append the transaction to be rolled back in case some of the ingresses were created
		handler.transactions = append(handler.transactions, txn)
	}
}

func (ch *CommandHandler) Deploy(command *cobra.Command, args []string) {

	name := args[0]
	image := args[1]
	configFile, _ := command.Flags().GetString("endpoint-config")
	port, _ := command.Flags().GetInt32("port")
	min, _ := command.Flags().GetInt32("min")
	max, _ := command.Flags().GetInt32("max")
	if max < min {
		max = min
	}

	fmt.Println("Starting deploy...")
	handler := NewDeployHandler(ch.Client, name)
	handler.createService(name, port)
	handler.createDeployment(name, image, port, min)
	handler.createAutoscaler(name, min, max)
	handler.createEndpoints(name, port, configFile)

	if handler.err != nil {
		if errors.IsAlreadyExists(handler.err) {
			fmt.Printf("\"%v\" already exists\n", name)
		} else {
			fmt.Printf("Unknown error: %v\n", handler.err)
		}
		fmt.Println("Rolling back...")
		handler.Rollback()
		fmt.Println("Rollback complete")
	} else {
		fmt.Println("Deploy complete")
	}
}