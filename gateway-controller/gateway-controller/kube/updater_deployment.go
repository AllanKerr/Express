package kube

import (
	typedappsv1beta2 "k8s.io/client-go/kubernetes/typed/apps/v1beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

type ContainerUpdate struct {
	Image *string
}

type DeploymentUpdater struct {
	dInterface typedappsv1beta2.DeploymentInterface
}

func NewDeploymentUpdater(client *Client, namespace string) *DeploymentUpdater {
	return &DeploymentUpdater{
		client.AppsV1beta2().Deployments(namespace),
	}
}

func (updater *DeploymentUpdater) GetModifiers() []string {
	return []string{
		"image",
	}
}

func (updater *DeploymentUpdater) Update(name string, update interface{}) error {

	containerUpdate := update.(*ContainerUpdate)
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
