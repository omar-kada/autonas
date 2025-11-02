package testutil

import (
	"reflect"
	"testing"
)

type Mocker struct {
	Calls [][]any
}

func (m *Mocker) Mock(method string, fn func()) {
	m.Calls = append(m.Calls, []any{method, fn})
}

func (m *Mocker) MockErr(method string, fn func()) {
	m.Calls = append(m.Calls, []any{method, fn})
	fn()
}

func (m *Mocker) AddCall(method string, args ...any) {
	call := []any{method}
	call = append(call, args...)
	m.Calls = append(m.Calls, call)
}

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
