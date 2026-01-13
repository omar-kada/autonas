package shell

import (
	"errors"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunCommand_Success(t *testing.T) {
	originalExecCommand := execCommand
	defer func() { execCommand = originalExecCommand }()

	execCommand = func(_ string, _ ...string) *exec.Cmd {
		return exec.Command("echo", "success")
	}

	err := NewExecutor().Exec("go", "help")
	assert.NoError(t, err)
}

func TestRunCommand_NoArgs(t *testing.T) {
	originalExecCommand := execCommand
	defer func() { execCommand = originalExecCommand }()

	execCommand = func(_ string, _ ...string) *exec.Cmd {
		c := exec.Command("false")
		c.Stderr = nil
		return c
	}

	err := NewExecutor().Exec("go")
	assert.ErrorContains(t, err, "exit status 1")
}

func TestRunCommand_NotFound(t *testing.T) {
	originalExecCommand := execCommand
	defer func() { execCommand = originalExecCommand }()

	execCommand = func(_ string, _ ...string) *exec.Cmd {
		return exec.Command("non-existent-command")
	}

	err := NewExecutor().Exec("dummyCmd")
	assert.ErrorContains(t, err, "executable not found")
}

func TestRunCommand_ExecError(t *testing.T) {
	originalExecCommand := execCommand
	defer func() { execCommand = originalExecCommand }()

	execCommand = func(_ string, _ ...string) *exec.Cmd {
		return &exec.Cmd{
			Path: "invalid-path",
			Err:  errors.New("exec error"),
		}
	}

	err := NewExecutor().Exec("echo")
	assert.ErrorContains(t, err, "exec error")
}
