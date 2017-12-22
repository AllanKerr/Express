package kube

import (
	apiv1 "k8s.io/api/core/v1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	typedappsv1beta2 "k8s.io/client-go/kubernetes/typed/apps/v1beta2"
	typedapiv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	autoscalingv2beta1 "k8s.io/api/autoscaling/v2beta1"
	typedautoscalingv2beta1 "k8s.io/client-go/kubernetes/typed/autoscaling/v2beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Transaction interface {
	Execute(interface{}) error
	Rollback(name string) error
}

type DeploymentTransaction struct {
	dInterface typedappsv1beta2.DeploymentInterface
}

func NewDeploymentTransaction(client *Client, namespace string) *DeploymentTransaction {
	return &DeploymentTransaction{
		client.AppsV1beta2().Deployments(namespace),
	}
}

func (txn *DeploymentTransaction) Execute(deployment interface{}) error {
	_, err := txn.dInterface.Create(deployment.(*appsv1beta2.Deployment))
	return err
}

func (txn *DeploymentTransaction) Rollback(name string) error {
	deletePolicy := metav1.DeletePropagationForeground
	return txn.dInterface.Delete(name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
}

type ServiceTransaction struct {
	sInterface typedapiv1.ServiceInterface
}

func NewServiceTransaction(client *Client, namespace string) *ServiceTransaction {
	return &ServiceTransaction{
		client.CoreV1().Services(namespace),
	}
}

func (txn *ServiceTransaction) Execute(service interface{}) error {
	_, err := txn.sInterface.Create(service.(*apiv1.Service))
	return err
}

func (txn *ServiceTransaction) Rollback(name string) error {
	deletePolicy := metav1.DeletePropagationForeground
	return txn.sInterface.Delete(name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
}

type AutoscalerTransaction struct {
	aInterface typedautoscalingv2beta1.HorizontalPodAutoscalerInterface
}

func NewAutoscalerTransaction(client *Client, namespace string) *AutoscalerTransaction {
	return &AutoscalerTransaction{
		client.AutoscalingV2beta1().HorizontalPodAutoscalers(namespace),
	}
}

func (txn *AutoscalerTransaction) Execute(service interface{}) error {
	_, err := txn.aInterface.Create(service.(*autoscalingv2beta1.HorizontalPodAutoscaler))
	return err
}

func (txn *AutoscalerTransaction) Rollback(name string) error {
	deletePolicy := metav1.DeletePropagationForeground
	return txn.aInterface.Delete(name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
}
