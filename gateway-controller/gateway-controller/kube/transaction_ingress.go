package kube

import (
	typedextensionsv1beta1 "k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type IngressTransaction struct {
	iInterface typedextensionsv1beta1.IngressInterface
}

func NewIngressTransaction(client *Client, namespace string) *IngressTransaction {
	return &IngressTransaction{
		iInterface: client.ExtensionsV1beta1().Ingresses(namespace),
	}
}

func (txn *IngressTransaction) Execute(ing interface{}) error {

	ingresses := ing.([]*extensionsv1beta1.Ingress)
	for _, ingress := range ingresses {
		if _, err := txn.iInterface.Create(ingress); err != nil {
			return err
		}
	}
	return nil
}

func (txn *IngressTransaction) Rollback(name string) error {

	deletePolicy := metav1.DeletePropagationForeground
	options := &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}

	return txn.iInterface.DeleteCollection(options, metav1.ListOptions{
		LabelSelector: "app=" + name,
	})
}
