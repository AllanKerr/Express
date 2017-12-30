package kube

import (
	"testing"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	typedappsv1beta2 "k8s.io/client-go/kubernetes/typed/apps/v1beta2"
)

type mockDeployment struct {
	typedappsv1beta2.DeploymentInterface
}

func (mock *mockDeployment) Create(deployment*appsv1beta2.Deployment) (*appsv1beta2.Deployment, error) {
	return deployment, nil
}

func TestDeploymentTransaction_Execute(t *testing.T) {

	// test error for unexpected execution type
	txn := DeploymentTransaction{
		&mockDeployment{},
	}
	err := txn.Execute("unexpected type")
	if err == nil {
		t.Error("Expected error for execution with unexpected type")
	}

	// test execute with a valid deployment
	err = txn.Execute(DefaultDeploymentConfig())
	if err != nil {
		t.Error("Error while executing autoscale transaction")
	}
}