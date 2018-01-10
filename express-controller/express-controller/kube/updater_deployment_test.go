package kube

import (
	"testing"
	typedappsv1beta2 "k8s.io/client-go/kubernetes/typed/apps/v1beta2"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/api/errors"
	apiv1 "k8s.io/api/core/v1"
)

// Test that the correct set of modifiers are being detected
func TestDeploymentModifiers(t *testing.T) {

	updater := DeploymentUpdater{}

	modifiers := map[string]int{
		"image": 0,
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

// Mock for simulating the deployment not being found
type mockDeploymentNotFound struct {
	typedappsv1beta2.DeploymentInterface
}

func (mock *mockDeploymentNotFound) Get(name string, options metav1.GetOptions) (*appsv1beta2.Deployment, error) {
	return nil, errors.NewNotFound(schema.GroupResource{}, "not found")
}

// Test updating a deployment that doesn't exist
// Expected not found error
func TestUpdateDeploymentNotFound(t *testing.T) {

	name := "testname"

	updater := DeploymentUpdater{
		&mockDeploymentNotFound{},
	}

	err := updater.Update(name, &ContainerUpdate{})
	if !errors.IsNotFound(err) {
		t.Errorf("Unexpected error for not found %v", err)
	}
}

// Test updating with the wrong update type
// Expected generic error without crash
func TestUpdateDeploymentWrongType(t *testing.T) {

	updater := DeploymentUpdater{}
	err := updater.Update("testname", "unexpected type")
	if err == nil {
		t.Errorf("Expected error for wrong type.")
	}
}

// Mock for finding the deployment and successful update
type mockDeploymentUpdaterSuccess struct {
	t *testing.T
	typedappsv1beta2.DeploymentInterface
}

func (mock *mockDeploymentUpdaterSuccess) Get(name string, options metav1.GetOptions) (*appsv1beta2.Deployment, error) {
	return &appsv1beta2.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: appsv1beta2.DeploymentSpec{
			Template: apiv1.PodTemplateSpec{
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{Image: "old"},
					},
				},
			},
		},
	}, nil
}

func (mock *mockDeploymentUpdaterSuccess) Update(autoscaler *appsv1beta2.Deployment) (*appsv1beta2.Deployment, error) {

	image := autoscaler.Spec.Template.Spec.Containers[0].Image
	if image != "new" {
		mock.t.Error("Unexpected deployment image %v", image)
	}
	name := autoscaler.ObjectMeta.Name
	if name != "testname" {
		mock.t.Error("Unexpected deployment name %v", name)
	}
	return autoscaler, nil
}

// Test successful update
func TestUpdateDeploymentSuccess(t *testing.T) {

	image := "new"
	update := &ContainerUpdate{
		Image: &image,
	}

	updater := DeploymentUpdater{
		&mockDeploymentUpdaterSuccess{t, nil},
	}

	err := updater.Update("testname", update)
	if err != nil {
		t.Errorf("Unexpected update error %v", err)
	}
}
