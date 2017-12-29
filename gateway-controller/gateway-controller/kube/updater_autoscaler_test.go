package kube

import "testing"

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

func TestUpdateSuccess(t *testing.T) {

	updater := AutoscalerUpdater{}
}
