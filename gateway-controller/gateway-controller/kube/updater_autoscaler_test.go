package kube

import (
	"testing"
	typedautoscalingv2beta1 "k8s.io/client-go/kubernetes/typed/autoscaling/v2beta1"
	v2beta1 "k8s.io/api/autoscaling/v2beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	testify "github.com/stretchr/testify/mock"
)

func TestModifiers(t *testing.T) {

	updater := AutoscalerUpdater{}

	modifiers := map[string]int{
		"min": 0,
		"max": 0,
	}

	for _, modifier := range updater.GetModifiers() {
		if _, ok := modifiers[modifier]; !ok {
			t.Error("Unexpected autoscaler modifier.")
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

// Mock for simulating the autoscaler not being found
type mockHorizontalPodAutoscalerNotFound struct {
	testify.Mock
	typedautoscalingv2beta1.HorizontalPodAutoscalerInterface
}

func (mock *mockHorizontalPodAutoscalerNotFound) Get(name string, options v1.GetOptions) (*v2beta1.HorizontalPodAutoscaler, error) {
	mock.Called(name)
	return nil, errors.NewNotFound(schema.GroupResource{}, "not found")
}

// Test updating an autoscaler that doesn't exist
// Expected not found error
func TestUpdateNotFound(t *testing.T) {

	name := "testname"

	mock := &mockHorizontalPodAutoscalerNotFound{}
	mock.On("Get", name).
		Return(nil, errors.NewNotFound(schema.GroupResource{}, "not found"))
	updater := AutoscalerUpdater{mock}

	err := updater.Update(name, &AutoscalerUpdate{})
	if !errors.IsNotFound(err) {
		t.Errorf("Unexpected error for not found %v", err)
	}
	mock.AssertCalled(t, "Get", "testname")
}

// Test updating with the wrong update type
// Expected generic error without crash
func TestUpdateWrongType(t *testing.T) {

	updater := AutoscalerUpdater{}
	err := updater.Update("testname", "unexpected type")
	if err == nil {
		t.Errorf("Expected error for wrong type.")
	}
}

// Mock for finding the autoscaler and successful update
type mockHorizontalPodAutoscalerSuccess struct {
	t *testing.T
	typedautoscalingv2beta1.HorizontalPodAutoscalerInterface
}

func (mock *mockHorizontalPodAutoscalerSuccess) Get(name string, options v1.GetOptions) (*v2beta1.HorizontalPodAutoscaler, error) {

	var min int32
	min = 1

	autoscaler := &v2beta1.HorizontalPodAutoscaler{}
	autoscaler.Name = "autoscaler"
	autoscaler.Spec = v2beta1.HorizontalPodAutoscalerSpec{}
	autoscaler.Spec.MinReplicas = &min
	autoscaler.Spec.MaxReplicas = 4
	return autoscaler, nil
}

func (mock *mockHorizontalPodAutoscalerSuccess) Update(autoscaler *v2beta1.HorizontalPodAutoscaler) (*v2beta1.HorizontalPodAutoscaler, error) {

	if autoscaler.Name != "autoscaler" {
		mock.t.Error("Unexpected name %v", autoscaler.Name)
	}
	if *autoscaler.Spec.MinReplicas != 2 {
		mock.t.Error("Unexpected minimum replicas %v", *autoscaler.Spec.MinReplicas)
	}
	if autoscaler.Spec.MaxReplicas != 4 {
		mock.t.Error("Unexpected maximum replicas %v", *autoscaler.Spec.MinReplicas)
	}
	return autoscaler, nil
}

// Test successful update
func TestUpdateSuccess(t *testing.T) {

	var min int32
	min = 2
	update := &AutoscalerUpdate{
		MinReplicas: &min,
		MaxReplicas: nil,
	}

	updater := AutoscalerUpdater{
		&mockHorizontalPodAutoscalerSuccess{t, nil},
	}

	err := updater.Update("testname", update)
	if err != nil {
		t.Errorf("Unexpected update error %v", err)
	}
}






