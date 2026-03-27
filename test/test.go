package test

import (
	"testing"
)

func Compare[T comparable](t testing.TB, variableName string, received, expected T) {
	compare(t, "", variableName, received, expected, func(received, expected T) bool {
		return received == expected
	})
}

func ComparePointers[T comparable](t testing.TB, variableName string, received, expected *T) {
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

func CompareErrors(t testing.TB, variableName string, received, expected error) {
	compare(t, "error", variableName, received, expected, func(received, expected error) bool {
		return received == expected
	})
}

func CompareCustomErrors(t testing.TB, variableName string, received, expected error) {
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

func isSlicesEqual[T, K any](received []T, expected []K, isEqual func(T, K) bool) bool {
	if (received == nil) != (expected == nil) {
		return false
	}
	if len(received) != len(expected) {
		return false
	}
	for i := range received {
		if !isEqual(received[i], expected[i]) {
			return false
		}
	}
	return true
}

func CompareSlicesWithFunc[T, K any](t testing.TB, variableName string, received []T, expected []K, isEqual func(T, K) bool) {
	if isEqual == nil {
		t.Fatal("isEqual can't be nil")
	}
	if isSlicesEqual(received, expected, isEqual) {
		return
	}
	t.Errorf("unexpected slice %s: got '%#v', expected '%#v'\n", variableName, received, expected)
}

func CompareSlices[T comparable](t testing.TB, variableName string, received, expected []T) {
	CompareSlicesWithFunc(t, variableName, received, expected, func(receivedItem, expectedItem T) bool {
		return receivedItem == expectedItem
	})
}

func compare[T any](
	t testing.TB,
	kind, variableName string,
	received, expected T,
	compare func(T, T) bool) {
	if ok := compare(received, expected); !ok {
		t.Errorf("unexpected %s %s: got '%#v', expected '%#v'\n", kind, variableName, received, expected)
	}
}
