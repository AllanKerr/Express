package kube

import (
	autoscalingv2beta1 "k8s.io/api/autoscaling/v2beta1"
	typedautoscalingv2beta1 "k8s.io/client-go/kubernetes/typed/autoscaling/v2beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"errors"
)

// Transaction to create a Kubernetes autoscaler
type AutoscalerTransaction struct {
	aInterface typedautoscalingv2beta1.HorizontalPodAutoscalerInterface
}

// Create a new transaction that operates on the specified namespace
func NewAutoscalerTransaction(client Client, namespace string) *AutoscalerTransaction {
	return &AutoscalerTransaction{
		client.AutoscalingV2beta1().HorizontalPodAutoscalers(namespace),
	}
}

// Execute the transaction by creating the specified object
func (txn *AutoscalerTransaction) Execute(object interface{}) error {

	autoscaler, ok := object.(*autoscalingv2beta1.HorizontalPodAutoscaler)
	if !ok {
		return errors.New("unexpected autoscaler execute type, expected *HorizontalPodAutoscaler")
	}
	_, err := txn.aInterface.Create(autoscaler)
	return err
}

// Rollback the transaction by deleting the a previously created autoscaler
// A transaction may be recreated then rolled back at a different point in time
// than when it was first created
func (txn *AutoscalerTransaction) Rollback(name string) error {
	deletePolicy := metav1.DeletePropagationForeground
	return txn.aInterface.Delete(name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
}
