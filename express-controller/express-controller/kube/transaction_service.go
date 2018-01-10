package kube

import (
	apiv1 "k8s.io/api/core/v1"
	typedapiv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"errors"
)

// Transaction to create a Kubernetes service
type ServiceTransaction struct {
	sInterface typedapiv1.ServiceInterface
}

// Create a new transaction that operates on the specified namespace
func NewServiceTransaction(client Client, namespace string) *ServiceTransaction {
	return &ServiceTransaction{
		client.CoreV1().Services(namespace),
	}
}

// Execute the transaction by creating a deployment from the provided object config
// The config must be a pointer to a Kubernetes Service
func (txn *ServiceTransaction) Execute(object interface{}) error {

	service, ok := object.(*apiv1.Service)
	if !ok {
		return errors.New("unexpected service execute object type, expected *Service")
	}
	_, err := txn.sInterface.Create(service)
	return err
}

// Rollback the transaction by deleting the a previously created service
// A transaction may be recreated then rolled back at a different point in time
// than when it was first created
func (txn *ServiceTransaction) Rollback(name string) error {
	deletePolicy := metav1.DeletePropagationForeground
	return txn.sInterface.Delete(name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
}
