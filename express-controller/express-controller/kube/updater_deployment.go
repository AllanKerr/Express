package kube

import (
	typedappsv1beta2 "k8s.io/client-go/kubernetes/typed/apps/v1beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
	"errors"
)

// The updater object for updating a deployment
// Only the image may be updated
// If image is nil then the old value is kept and the deployment is not updated.
type ContainerUpdate struct {
	Image *string
}

// The deployment update handler
type DeploymentUpdater struct {
	dInterface typedappsv1beta2.DeploymentInterface
}

// Create a new deployment updater to update a deployment on the specified namespace
func NewDeploymentUpdater(client Client, namespace string) *DeploymentUpdater {
	return &DeploymentUpdater{
		client.AppsV1beta2().Deployments(namespace),
	}
}

// The list of modifiers that will result in a deployment update
func (updater *DeploymentUpdater) GetModifiers() []string {
	return []string{
		"image",
	}
}

// Performs an update on the existing deployment with the specified name
// The update value must be a *ContainerUpdate
func (updater *DeploymentUpdater) Update(name string, update interface{}) error {

	containerUpdate, ok := update.(*ContainerUpdate)
	if !ok {
		return errors.New("unexpected update type, expected *ContainerUpdate")
	}
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {

		result, getErr := updater.dInterface.Get(name, metav1.GetOptions{})
		if getErr != nil {
			return getErr
		}

		containerSpec := result.Spec.Template.Spec.Containers[0]
		if containerUpdate.Image != nil {
			containerSpec.Image = *containerUpdate.Image
		}
		result.Spec.Template.Spec.Containers[0] = containerSpec

		_, updateErr := updater.dInterface.Update(result)
		return updateErr
	})
}
