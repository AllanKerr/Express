package kube

import (
	"testing"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strconv"
)

// Test set port to check that all ports for an Ingress configuration are changed
func TestIngressUpdater_setPort(t *testing.T) {

	ingress := DefaultIngressConfig()

	var paths []v1beta1.HTTPIngressPath
	for i := 0; i < 5; i++ {
		paths = append(paths, v1beta1.HTTPIngressPath{
			Path: "test/path" + strconv.Itoa(i),
			Backend: v1beta1.IngressBackend{
				ServicePort: intstr.FromInt(80),
				ServiceName: "service",
			},
		})
	}
	ingress.Spec.Rules[0].HTTP.Paths = paths

	setPort(ingress, 8080)

	for i := 0; i < 5; i++ {
		path := ingress.Spec.Rules[0].HTTP.Paths[i]

		if path.Backend.ServicePort.IntVal != 8080 {
			t.Error("Error setting port")
		}
		if path.Backend.ServiceName != "service" {
			t.Errorf("Error, unexpected service name: %v", path.Backend.ServiceName)
		}
		if path.Path != "test/path" + strconv.Itoa(i) {
			t.Errorf("Error, unexpected service path: %v", path.Path)
		}
	}
}

func TestIngressUpdater_findAndRemove(t *testing.T) {

	var ingresses[]v1beta1.Ingress
	for i := 0; i < 5; i++ {
		config := DefaultIngressConfig()
		config.Labels = map[string] string {
			"identifier" : strconv.Itoa(i),
		}
		ingresses = append(ingresses, *config)
	}


	// Test finding an Ingress not in the slice
	notFoundConfig := DefaultIngressConfig()
	notFoundConfig.Labels = map[string] string {
		"identifier" : strconv.Itoa(-1),
	}
	found, newIngresses := findAndRemove(notFoundConfig, ingresses)
	if found {
		t.Error("Error, found and removed Ingress that shouldn't have been found")
	}
	if len(newIngresses) != len(ingresses) {
		t.Error("Error, Ingress unexpectedly removed")
	}

	// Test finding an Ingress in the slice
	foundConfig := DefaultIngressConfig()
	foundConfig.Labels = map[string] string {
		"identifier" : strconv.Itoa(2),
	}
	found, newIngresses = findAndRemove(foundConfig, ingresses)
	if !found {
		t.Error("Error, found and removed Ingress that shouldn't have been found")
	}
	if len(newIngresses) != len(ingresses) - 1 {
		t.Errorf("Error, Ingress unexpected length: %v", len(newIngresses))
	}
}

// Test that the correct set of modifiers are being detected
func TestIngressUpdater_GetModifiers(t *testing.T) {

	updater := IngressUpdater{}

	modifiers := map[string]int{
		"endpoint-config": 0,
	}

	for _, modifier := range updater.GetModifiers() {
		if _, ok := modifiers[modifier]; !ok {
			t.Error("Unexpected deployment modifier.")
		} else {
			modifiers[modifier]++
		}
	}

	for modifier, n := range modifiers {
		if n == 0 {
			t.Errorf("Modifier not found %v", modifier)
		} else if n > 1 {
			t.Errorf("Duplicate modifier found %v", modifier)
		}
	}
}
