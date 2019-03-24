package main

import (
	"testing"
)

func TestValidateInvite(t *testing.T) {

	expected := false
	input := "hello world"
	result := validateInvite(input)
	if result != expected {
		t.Errorf("Expected %q to produce status %t, got %t", input, expected, result)
		return
	}

	expected = true
	input = "xzbcdfghj"
	result = validateInvite(input)
	if result != expected {
		t.Errorf("Expected %q to produce status %t, got %t", input, expected, result)
		return
	}

	expected = true
	input = "xzb-cdf-ghj"
	result = validateInvite(input)
	if result != expected {
		t.Errorf("Expected %q to produce status %t, got %t", input, expected, result)
		return
	}
}
