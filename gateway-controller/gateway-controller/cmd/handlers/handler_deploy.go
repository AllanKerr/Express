package handlers

import (
	"github.com/spf13/cobra"
	"fmt"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	autoscalingv2beta1 "k8s.io/api/autoscaling/v2beta1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"gateway-controller/kube"
	"k8s.io/apimachinery/pkg/api/errors"
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

	labels := map[string]string{
		"app": name,
	}

	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: labels,
		},
		Spec: apiv1.ServiceSpec{
			Ports: []apiv1.ServicePort{
				{
					Protocol: apiv1.ProtocolTCP,
					Port: port,
					TargetPort: intstr.FromInt(int(port)),
				},
			},
			Selector: labels,
			Type: apiv1.ServiceTypeNodePort,
		},
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

	labels := map[string]string{
		"app": name,
	}

	autoscaler := &autoscalingv2beta1.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: labels,
		},
		Spec: autoscalingv2beta1.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv2beta1.CrossVersionObjectReference{
				Kind: "Deployment",
				Name: name,
			},
			MinReplicas: &min,
			MaxReplicas: max,
		},
	}
	txn := kube.NewAutoscalerTransaction(handler.client, apiv1.NamespaceDefault)
	if err := handler.Execute(txn, autoscaler); err == nil {
		fmt.Printf("Created autoscaler %q.\n", name)
	}
}

func (ch *CommandHandler) Deploy(command *cobra.Command, args []string) {

	name := args[0]
	image := args[1]
	port, _ := command.Flags().GetInt32("port")
	min, _ := command.Flags().GetInt32("min")
	max, _ := command.Flags().GetInt32("max")
	if max < min {
		max = min
	}

	handler := NewDeployHandler(ch.Client, name)
	handler.createService(name, port)
	handler.createDeployment(name, image, port, min)
	handler.createAutoscaler(name, min, max)

	if handler.err != nil {
		if errors.IsAlreadyExists(handler.err) {
			fmt.Printf("\"%v\" already exists\n", name)
		} else {
			fmt.Printf("Unknown error: %v\n", handler.err)
		}
		fmt.Println("Rolling back...")
		handler.Rollback()
	}
}