package main

import (
	"errors"
	"testing"
)

type TestUiReader struct {
	mockReturn string
	mockErr    error
}

func (t TestUiReader) ReadString(b byte) (s string, err error) {
	if t.mockReturn != "" {
		return t.mockReturn, nil
	}
	return "", t.mockErr
}

func TestUserInput_Ask(t *testing.T) {

	mock := TestUiReader{"hello, world\n", nil}
	in := UserInput{mock}
	expected := "hello,world"
	result, err := in.Ask("", 1, func(input string) (valid bool) {
		return true
	}, "")
	if err != nil {
		t.Errorf("Expected %q, got error", expected)
		return
	}
	if result != expected {
		t.Errorf("Expected %q, got %t", expected, result)
		return
	}
}

func TestUserInput_AskReturnsErr(t *testing.T) {

	errMock := errors.New("something went wrong")
	mock := TestUiReader{"", errMock}
	in := UserInput{mock}
	result, err := in.Ask("", 1, func(input string) (valid bool) {
		return false
	}, "")
	if err == nil {
		t.Errorf("Expected %q, got nil", errMock.Error())
		return
	}
	if result != "" {
		t.Errorf("Expected empty string, got %t", result)
		return
	}
}

func TestUserInput_AskYesOrNo(t *testing.T) {

	inputs := []string{"yes", "y", "Y", "YES", "no", "N", "n", "No"}
	outputs := []bool{true, true, true, true, false, false, false, false}

	for i := 0; i < len(inputs); i++ {
		mock := TestUiReader{inputs[i], nil}
		in := UserInput{mock}
		result, err := in.AskYesOrNo("", 1)

		if err != nil {
			t.Errorf("Expected nil error, got %q", err.Error())
		}

		if result != outputs[i] {
			t.Errorf("Expected %t, got %t", result, outputs[i])
		}
	}

}
