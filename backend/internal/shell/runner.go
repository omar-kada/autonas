// Package shell provides utilities to interact with shell.
package shell

import (
	"fmt"
	"log/slog"
	"os/exec"

	"omar-kada/autonas/internal/events"
)

// Executor abstracts writing content to a file
type Executor interface {
	Exec(cmd string, args ...string) error
}

type cmdExecuter struct{}

// NewExecutor creates and new Writer and returns it
func NewExecutor() Executor {
	return cmdExecuter{}
}

// Run runs a shell command and returns error if any
func (cmdExecuter) Exec(cmd string, args ...string) error {
	path, err := exec.LookPath(cmd)
	if err != nil {
		return fmt.Errorf("executable not found: %w", err)
	}
	c := execCommand(path, args...)
	c.Stdout = events.NewSlogWriter(slog.LevelInfo)
	c.Stderr = events.NewSlogWriter(slog.LevelError)
	return c.Run()
}

// execCommand is a wrapper for exec.Command for testability
var execCommand = defaultExecCommand

func defaultExecCommand(cmd string, args ...string) *exec.Cmd {
	return exec.Command(cmd, args...)
}
