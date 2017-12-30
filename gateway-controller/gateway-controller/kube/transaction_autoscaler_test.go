package kube

import (
	"testing"
	typedautoscalingv2beta1 "k8s.io/client-go/kubernetes/typed/autoscaling/v2beta1"
	autoscalingv2beta1 "k8s.io/api/autoscaling/v2beta1"
)

type mockHorizontalPodAutoscaler struct {
	typedautoscalingv2beta1.HorizontalPodAutoscalerInterface
}

func (mock *mockHorizontalPodAutoscaler) Create(autoscaler*autoscalingv2beta1.HorizontalPodAutoscaler) (*autoscalingv2beta1.HorizontalPodAutoscaler, error) {
	return autoscaler, nil
}

func TestAutoscalerTransaction_Execute(t *testing.T) {

	// test error for unexpected execution type
	txn := AutoscalerTransaction{
		&mockHorizontalPodAutoscaler{},
	}
	err := txn.Execute("unexpcted type")
	if err == nil {
		t.Error("Expected error for execution with unexpected type")
	}

	// test execute with a valid autoscaler
	err = txn.Execute(DefaultAutoscalerConfig())
	if err != nil {
		t.Error("Error while executing autoscale transaction")
	}
}