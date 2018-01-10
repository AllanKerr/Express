package kube

import (
	"testing"
	apiv1 "k8s.io/api/core/v1"
	typedapiv1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type mockService struct {
	typedapiv1.ServiceInterface
}

func (mock *mockService) Create(service*apiv1.Service) (*apiv1.Service, error) {
	return service, nil
}

func TestServiceTransaction_Execute(t *testing.T) {

	// test error for unexpected execution type
	txn := ServiceTransaction{
		&mockService{},
	}
	err := txn.Execute("unexpected type")
	if err == nil {
		t.Error("Expected error for execution with unexpected type")
	}

	// test execute with a valid autoscaler
	err = txn.Execute(DefaultServiceConfig())
	if err != nil {
		t.Error("Error while executing autoscale transaction")
	}
}