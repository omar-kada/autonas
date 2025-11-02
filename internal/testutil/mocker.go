package testutil

import (
	"reflect"
	"testing"
)

// Mocker records method calls and their arguments for testing purposes.
type Mocker struct {
	Calls [][]any
}

// Reset clears all recorded calls.
func (m *Mocker) Reset() {
	m.Calls = nil
}

// AddCall records a method call with its arguments.
func (m *Mocker) AddCall(method string, args ...any) {
	call := []any{method}
	call = append(call, args...)
	m.Calls = append(m.Calls, call)
}

// AssertCalls checks if the recorded calls match the expected calls.
func (m *Mocker) AssertCalls(t *testing.T, expected [][]any) {
	t.Helper()
	if len(m.Calls) != len(expected) {
		t.Errorf("expected %d calls, got %d", len(expected), len(m.Calls))
		return
	}
	for i, call := range m.Calls {
		if !reflect.DeepEqual(call, expected[i]) {
			t.Errorf("call %d: expected %v, got %v", i, expected[i], call)
		}
	}
}
