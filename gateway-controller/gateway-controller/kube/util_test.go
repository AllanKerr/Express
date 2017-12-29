package kube

import "testing"

func TestHashString(t *testing.T) {

	if HashString("") != HashString("") {
		t.Fatal("Hash codes for empty strings did not match.")
	}

	if HashString("a string") != HashString("a string") {
		t.Fatal("Hash codes for non-empty strings did not match.")
	}

	if HashString("differing") == HashString("differin") {
		t.Fatal("Hash codes for slightly different strings did not differ.")
	}

	if HashString("") == HashString("drastically") {
		t.Fatal("Hash codes for drastically different strings did not differ.")
	}
}
