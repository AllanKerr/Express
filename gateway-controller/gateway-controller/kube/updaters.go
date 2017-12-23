package kube

import (
	typedappsv1beta2 "k8s.io/client-go/kubernetes/typed/apps/v1beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

type ContainerUpdate struct {
	Image *string
	Ports []ContainerPortUpdate
}

type ContainerPortUpdate struct {
	ContainerPort *int32
}

type ObjectUpdater interface {
	GetModifiers() []string
	Update(name string, spec interface{}) error
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
		"port",
	}
}

func (updater *DeploymentUpdater) Update(name string, update interface{}) error {

	return retry.RetryOnConflict(retry.DefaultRetry, func() error {

		result, getErr := updater.dInterface.Get(name, metav1.GetOptions{})
		if getErr != nil {
			return getErr
		}

		containerUpdate := update.(*ContainerUpdate)
		containerSpec := result.Spec.Template.Spec.Containers[0]
		if containerUpdate.Image != nil {
			containerSpec.Image = *containerUpdate.Image
		}
		result.Spec.Template.Spec.Containers[0] = containerSpec

		portUpdate := containerUpdate.Ports[0]
		portSpec := containerSpec.Ports[0]
		if portUpdate.ContainerPort != nil {
			portSpec.ContainerPort = *portUpdate.ContainerPort
		}
		result.Spec.Template.Spec.Containers[0].Ports[0] = portSpec

		_, updateErr := updater.dInterface.Update(result)
		return updateErr
	})
}
