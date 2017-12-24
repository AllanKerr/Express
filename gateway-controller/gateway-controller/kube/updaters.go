package kube

import (
	typedautoscalingv2beta1 "k8s.io/client-go/kubernetes/typed/autoscaling/v2beta1"
	typedappsv1beta2 "k8s.io/client-go/kubernetes/typed/apps/v1beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	typedapiv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	typedextensionsv1beta1 "k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/client-go/util/retry"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/intstr"
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

type IngressUpdater struct {
	iInterface typedextensionsv1beta1.IngressInterface
	sInterface 	typedapiv1.ServiceInterface
}

func NewIngressUpdater(client *Client, namespace string) *IngressUpdater {
	return &IngressUpdater{
		client.ExtensionsV1beta1().Ingresses(namespace),
		client.CoreV1().Services(namespace),
	}
}

func setPort(ingress *extensionsv1beta1.Ingress, port int32) {
	for i, path := range ingress.Spec.Rules[0].HTTP.Paths {
		path.Backend.ServicePort = intstr.FromInt(int(port))
		ingress.Spec.Rules[0].HTTP.Paths[i] = path
	}
}

func findAndRemove(ingress *extensionsv1beta1.Ingress, ingresses []extensionsv1beta1.Ingress) (bool, []extensionsv1beta1.Ingress) {

	identifier := ingress.Labels["identifier"]
	for i, existing := range ingresses {
		if existing.Labels["identifier"] ==	identifier {
			ingress.ObjectMeta.Name = existing.ObjectMeta.Name
			ingresses[i] = ingresses[0]
			return true, ingresses[1:]
		}
	}
	return false, ingresses
}

func (updater *IngressUpdater) getPort(name string) (int32, error) {
	if service, err := updater.sInterface.Get(name, metav1.GetOptions{}); err != nil {
		return 0, err
	} else {
		return service.Spec.Ports[0].Port , nil
	}
}

func (updater *IngressUpdater) getIngresses(name string) ([]extensionsv1beta1.Ingress, error) {
	ingressList, err := updater.iInterface.List(metav1.ListOptions{
		LabelSelector: "app=" + name,
	})
	if err != nil {
		return nil, err
	}
	return ingressList.Items, nil
}

func (updater *IngressUpdater) deleteIngresses(ingresses []extensionsv1beta1.Ingress) error {

	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}
	var err error
	for _, ingress := range ingresses {
		err = updater.iInterface.Delete(ingress.GetObjectMeta().GetName(), deleteOptions)
	}
	return err
}

func (updater *IngressUpdater) GetModifiers() []string {
	return []string{
		"endpoint-config",
	}
}

func (updater *IngressUpdater) Update(name string, update interface{}) error {

	port, portErr := updater.getPort(name)
	if portErr != nil {
		return portErr
	}
	ingresses, ingressErr := updater.getIngresses(name)
	if ingressErr != nil {
		return ingressErr
	}

	var found bool
	ingressUpdates := update.([]*extensionsv1beta1.Ingress)
	for _, ingress := range ingressUpdates {

		setPort(ingress, port)
		found, ingresses = findAndRemove(ingress, ingresses)

		var err error
		if found {
			err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
				_, updateErr := updater.iInterface.Update(ingress)
				return updateErr
			})
		} else {
			 _, err = updater.iInterface.Create(ingress)
		}
		if err != nil {
			return err
		}
	}
	return updater.deleteIngresses(ingresses)
}
