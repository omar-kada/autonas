package testutil

import (
	"testing"
	"time"
)

// WaitFor waits for the predicate to return true, up to the specified timeout.
// It returns the result of the predicate and true if the predicate succeeds within the timeout.
// If the timeout is reached, it returns the zero value of T and false.
func WaitFor[T any](timeout time.Duration, predicate func() (T, bool)) (T, bool) {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		if res, ok := predicate(); ok {
			return res, true
		}
		time.Sleep(200 * time.Millisecond)
	}
	var zero T
	return zero, false
}

// WaitForChannel waits for a signal on the given channel, up to the specified timeout.
// It returns the value received on the channel and true if successful.
// If the timeout is reached, it fails the test with the provided message and returns the zero value with false.
func WaitForChannel[T any](t testing.TB, done <-chan T, timeout time.Duration, message string) (T, bool) {
	select {
	case val := <-done:
		return val, true
	case <-time.After(timeout):
		t.Fatal(message)
		var zero T
		return zero, false
	}
}
