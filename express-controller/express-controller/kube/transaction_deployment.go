package kube

import (
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	typedappsv1beta2 "k8s.io/client-go/kubernetes/typed/apps/v1beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"errors"
)

// Transaction to create a Kubernetes deployment
type DeploymentTransaction struct {
	dInterface typedappsv1beta2.DeploymentInterface
}

// Create a new transaction that operates on the specified namespace
func NewDeploymentTransaction(client Client, namespace string) *DeploymentTransaction {
	return &DeploymentTransaction{
		client.AppsV1beta2().Deployments(namespace),
	}
}

// Execute the transaction by creating a deployment from the provided object config
// The config must be a pointer to a Kubernetes Deployment
func (txn *DeploymentTransaction) Execute(object interface{}) error {

	deployment, ok := object.(*appsv1beta2.Deployment)
	if !ok {
		return errors.New("unexpected autoscaler execute type, expected *Deployment")
	}
	_, err := txn.dInterface.Create(deployment)
	return err
}

// Rollback the transaction by deleting the a previously created deployment
// A transaction may be recreated then rolled back at a different point in time
// than when it was first created
func (txn *DeploymentTransaction) Rollback(name string) error {
	deletePolicy := metav1.DeletePropagationForeground
	return txn.dInterface.Delete(name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
}