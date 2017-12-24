package kube

import (
	apiv1 "k8s.io/api/core/v1"
	typedapiv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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
