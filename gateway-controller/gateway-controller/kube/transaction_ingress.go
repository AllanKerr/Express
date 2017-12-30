package kube

import (
	typedextensionsv1beta1 "k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"errors"
)

// Transaction to create new Kubernetes Ingresses
type IngressTransaction struct {
	iInterface typedextensionsv1beta1.IngressInterface
}

// Create a new transaction that operates on the specified namespace
func NewIngressTransaction(client Client, namespace string) *IngressTransaction {
	return &IngressTransaction{
		iInterface: client.ExtensionsV1beta1().Ingresses(namespace),
	}
}

// Execute the transaction by creating a set of Ingresses from the provided object config
// The config must be an array of pointers to Kubernetes Ingress configurations
func (txn *IngressTransaction) Execute(ing interface{}) error {

	ingresses, ok := ing.([]*extensionsv1beta1.Ingress)
	if !ok {
		return errors.New("unexpected Ingress execute object type, expected []*Ingress")
	}
	// Create each Ingress
	for _, ingress := range ingresses {
		if _, err := txn.iInterface.Create(ingress); err != nil {
			// Fail if creating any of them fails
			return err
		}
	}
	return nil
}

// Rollback the transaction by deleting the a previously created ingresses
// A transaction may be recreated then rolled back at a different point in time
// than when it was first created
func (txn *IngressTransaction) Rollback(name string) error {

	deletePolicy := metav1.DeletePropagationForeground
	options := &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}

	// Delete all of the Ingresses for the deployment at once
	return txn.iInterface.DeleteCollection(options, metav1.ListOptions{
		LabelSelector: "app=" + name,
	})
}
