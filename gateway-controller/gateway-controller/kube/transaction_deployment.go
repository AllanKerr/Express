package kube

import (
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	typedappsv1beta2 "k8s.io/client-go/kubernetes/typed/apps/v1beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeploymentTransaction struct {
	dInterface typedappsv1beta2.DeploymentInterface
}

func NewDeploymentTransaction(client *Client, namespace string) *DeploymentTransaction {
	return &DeploymentTransaction{
		client.AppsV1beta2().Deployments(namespace),
	}
}

func (txn *DeploymentTransaction) Execute(deployment interface{}) error {
	_, err := txn.dInterface.Create(deployment.(*appsv1beta2.Deployment))
	return err
}

func (txn *DeploymentTransaction) Rollback(name string) error {
	deletePolicy := metav1.DeletePropagationForeground
	return txn.dInterface.Delete(name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
}