package kube

import (
	typedautoscalingv2beta1 "k8s.io/client-go/kubernetes/typed/autoscaling/v2beta1"
	typedappsv1beta2 "k8s.io/client-go/kubernetes/typed/apps/v1beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
	"github.com/pkg/errors"
)

type ContainerUpdate struct {
	Image *string
}

type AutoscalerUpdate struct {
	MinReplicas *int32
	MaxReplicas *int32
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

type AutoscalerUpdater struct {
	aInterface typedautoscalingv2beta1.HorizontalPodAutoscalerInterface
}

func NewAutoscalerUpdater(client *Client, namespace string) *AutoscalerUpdater {
	return &AutoscalerUpdater{
		client.AutoscalingV2beta1().HorizontalPodAutoscalers(namespace),
	}
}

func (updater *AutoscalerUpdater) GetModifiers() []string {
	return []string {
		"min",
		"max",
	}
}

func (updater *AutoscalerUpdater) Update(name string, update interface{}) error {

	autoscalerUpdate := update.(*AutoscalerUpdate)
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {

		result, getErr := updater.aInterface.Get(name, metav1.GetOptions{})
		if getErr != nil {
			return getErr
		}

		if autoscalerUpdate.MinReplicas != nil {
			result.Spec.MinReplicas = autoscalerUpdate.MinReplicas
		}
		if autoscalerUpdate.MaxReplicas != nil {
			result.Spec.MaxReplicas = *autoscalerUpdate.MaxReplicas
		}
		if result.Spec.MinReplicas != nil && (result.Spec.MaxReplicas < *result.Spec.MinReplicas) {
			return errors.New("max replicas must be greater than min replicas")
		}
		_, updateErr := updater.aInterface.Update(result)
		return updateErr
	})
	return nil
}
