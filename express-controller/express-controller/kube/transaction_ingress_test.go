package kube

import (
	typedextensionsv1beta1 "k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"testing"
)

type mockIngress struct {
	nCreated int
	typedextensionsv1beta1.IngressInterface
}

func (mock *mockIngress) Create(ingress *extensionsv1beta1.Ingress) (*extensionsv1beta1.Ingress, error) {
	mock.nCreated++
	return ingress, nil
}

func TestIngressTransaction_Execute(t *testing.T) {

	// test error for unexpected execution type
	mock := &mockIngress{}
	txn := IngressTransaction{
		mock,
	}
	err := txn.Execute("unexpected type")
	if err == nil {
		t.Error("Expected error for execution with unexpected type")
	}

	var ingresses []*extensionsv1beta1.Ingress
	for i := 0; i < 5; i++ {
		ingresses = append(ingresses, DefaultIngressConfig())
	}

	// test execute with a valid autoscaler
	err = txn.Execute(ingresses)
	if err != nil {
		t.Error("Error while executing autoscale transaction")
	}
	if mock.nCreated != 5 {
		t.Errorf("Error created unexpected number of transactions %v", mock.nCreated)
	}
}