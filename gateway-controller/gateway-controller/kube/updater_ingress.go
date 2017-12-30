package kube

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	typedapiv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	typedextensionsv1beta1 "k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/client-go/util/retry"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// Updater to update the existing set of Ingress configurations for a deployed application container
type IngressUpdater struct {
	iInterface typedextensionsv1beta1.IngressInterface
	sInterface 	typedapiv1.ServiceInterface
}

// Create a new updater to update the existing set of Ingress configurations in the specified namespace
// for a deployed application container
func NewIngressUpdater(client *Client, namespace string) *IngressUpdater {
	return &IngressUpdater{
		client.ExtensionsV1beta1().Ingresses(namespace),
		client.CoreV1().Services(namespace),
	}
}

// Set the service port for all paths in an Ingress configuration
func setPort(ingress *extensionsv1beta1.Ingress, port int32) {
	for i, path := range ingress.Spec.Rules[0].HTTP.Paths {
		path.Backend.ServicePort = intstr.FromInt(int(port))
		ingress.Spec.Rules[0].HTTP.Paths[i] = path
	}
}

// Find an Ingress configuration within a slice of configurations based on its identifier label
// If the Ingress configuration is found then it is removed from the slice
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

// Determines the port that should be routed to by getting
// the port for the deployed application container's service
func (updater *IngressUpdater) getPort(name string) (int32, error) {
	if service, err := updater.sInterface.Get(name, metav1.GetOptions{}); err != nil {
		return 0, err
	} else {
		return service.Spec.Ports[0].Port , nil
	}
}

// Get the set of Ingress configurations that exist for a deployed application
func (updater *IngressUpdater) getIngresses(name string) ([]extensionsv1beta1.Ingress, error) {
	ingressList, err := updater.iInterface.List(metav1.ListOptions{
		LabelSelector: "app=" + name,
	})
	if err != nil {
		return nil, err
	}
	return ingressList.Items, nil
}

// Delete the slice of Ingress configurations
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

// The list of modifiers that will result in an Ingress update
func (updater *IngressUpdater) GetModifiers() []string {
	return []string{
		"endpoint-config",
	}
}

// Performs an update on the existing Ingress configuration with the specified name
// The update value must be a []*Ingress for the set of Ingresses to be updated to.
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
			// update the Ingress configuration if there is already an Ingress
			// configuration with a matching identifier
			err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
				_, updateErr := updater.iInterface.Update(ingress)
				return updateErr
			})
		} else {
			// there is no existing Ingress configuration for the identifier so create one
			_, err = updater.iInterface.Create(ingress)
		}
		if err != nil {
			return err
		}
	}
	// delete any existing Ingress configurations that did not exist in the Ingress update set
	return updater.deleteIngresses(ingresses)
}
