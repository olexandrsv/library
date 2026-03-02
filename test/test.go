package test

import (
	"testing"
)

func Compare[T comparable](t *testing.T, variableName string, received, expected T) {
	compare(t, "", variableName, received, expected, func(received, expected T) bool {
		return received == expected
	})
}

func ComparePointers[T comparable](t *testing.T, variableName string, received, expected *T) {
	compare(t, "", variableName, received, expected, func(received, expected *T) bool {
		if received == expected {
			return true
		}
		if received == nil || expected == nil {
			return false
		}
		return *received == *expected
	})
}

func CompareErrors(t *testing.T, variableName string, received, expected error) {
	compare(t, "error", variableName, received, expected, func(received, expected error) bool {
		return received == expected
	})
}

func CompareCustomErrors(t *testing.T, variableName string, received, expected error) {
	compare(t, "error", variableName, received, expected, func(received, expected error) bool {
		if expected == received {
			return true
		}
		if expected == nil || received == nil {
			return false
		}
		return expected.Error() == received.Error()
	})
}

func CompareSlices[T comparable](t *testing.T, variableName string, received, expected []T) {
	compare(t, "slice", variableName, received, expected, func(received, expected []T) bool {
		if len(received) != len(expected) {
			return false
		}
		for _, v := range expected {
			for _, v2 := range received {
				if v != v2 {
					return false
				}
			}
		}
		return true
	})
}

func compare[T any](
	t *testing.T,
	kind, variableName string,
	received, expected T,
	compare func(T, T) bool) {
	if ok := compare(received, expected); !ok {
		t.Errorf("unexpected %s %s: got '%+v', expected '%+v'\n", kind, variableName, received, expected)
	}
}
