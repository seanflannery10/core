package assert

import (
	"reflect"
	"strings"
	"testing"
)

func Equal[T comparable](t *testing.T, actual, expected T) bool {
	t.Helper()

	if actual != expected {
		t.Errorf("not equal:\n"+
			"  actual: %v\n"+
			"expected: %v", actual, expected)

		return false
	}

	return true
}

func NotEqual[T comparable](t *testing.T, actual, expected T) bool {
	t.Helper()

	if actual == expected {
		t.Errorf("got for both: %v", actual)

		return false
	}

	return true
}

func SameType[T comparable](t *testing.T, actual, expected T) bool {
	t.Helper()

	if reflect.TypeOf(actual) != reflect.TypeOf(expected) {
		t.Errorf("not the same type:\n"+
			"  actual: %v\n"+
			"expected: %v", actual, expected)

		return false
	}

	return true
}

func Contains(t *testing.T, actual, expectedSubstring string) bool {
	t.Helper()

	if !strings.Contains(actual, expectedSubstring) {
		t.Errorf("%q does not contain %q", actual, expectedSubstring)
		return false
	}

	return true
}

func NotContains(t *testing.T, actual, expectedSubstring string) bool {
	t.Helper()

	if strings.Contains(actual, expectedSubstring) {
		t.Errorf("%q should not contain %q", actual, expectedSubstring)
		return false
	}

	return true
}

func NilError(t *testing.T, actual error) bool {
	t.Helper()

	if actual != nil {
		t.Errorf("%v should be nil", actual)

		return false
	}

	return true
}
