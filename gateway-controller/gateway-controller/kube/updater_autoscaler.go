package kube

import (
	typedautoscalingv2beta1 "k8s.io/client-go/kubernetes/typed/autoscaling/v2beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
	"github.com/pkg/errors"
)

// Autoscaler update used to update an existing autoscaler
// Only the min and max replicas can be updated
// If either value is nil then that means the value should not be updated
type AutoscalerUpdate struct {
	MinReplicas *int32
	MaxReplicas *int32
}

// The autoscaler update handler
type AutoscalerUpdater struct {
	aInterface typedautoscalingv2beta1.HorizontalPodAutoscalerInterface
}

// Create a new autoscaler updater to update an autoscaler on the specified namespace
func NewAutoscalerUpdater(client Client, namespace string) *AutoscalerUpdater {
	return &AutoscalerUpdater{
		client.AutoscalingV2beta1().HorizontalPodAutoscalers(namespace),
	}
}

// The list of modifiers that will result in an autoscaler update
func (updater *AutoscalerUpdater) GetModifiers() []string {
	return []string {
		"min",
		"max",
	}
}

// Performs an update on the autoscaler with the specified name
// The update value must be an *AutoscalerUpdate
func (updater *AutoscalerUpdater) Update(name string, update interface{}) error {

	autoscalerUpdate, ok := update.(*AutoscalerUpdate)
	if !ok {
		return errors.New("Unexpected update type, expected AutoscalerUpdate")
	}
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {

		result, getErr := updater.aInterface.Get(name, metav1.GetOptions{})
		if getErr != nil {
			return getErr
		}

		// Only update the values that were non-nil in the autoscaler update
		if autoscalerUpdate.MinReplicas != nil {
			result.Spec.MinReplicas = autoscalerUpdate.MinReplicas
		}
		if autoscalerUpdate.MaxReplicas != nil {
			result.Spec.MaxReplicas = *autoscalerUpdate.MaxReplicas
		}

		// Check that min < max
		if result.Spec.MinReplicas != nil && (result.Spec.MaxReplicas < *result.Spec.MinReplicas) {
			return errors.New("max replicas must be greater than min replicas")
		}
		_, updateErr := updater.aInterface.Update(result)
		return updateErr
	})
	return nil
}
