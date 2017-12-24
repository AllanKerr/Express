package kube

import (
	autoscalingv2beta1 "k8s.io/api/autoscaling/v2beta1"
	typedautoscalingv2beta1 "k8s.io/client-go/kubernetes/typed/autoscaling/v2beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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
