// Package shell provides utilities to interact with shell.
package shell

import (
	"fmt"
	"os"
	"os/exec"
)

// RunCommand runs a shell command and returns error if any
func RunCommand(cmdAndArgs ...string) error {
	c, err := execCommand(cmdAndArgs[0], cmdAndArgs[1:]...)
	if err != nil {
		return err
	}
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

// execCommand is a wrapper for exec.Command for testability
var execCommand = defaultExecCommand

func defaultExecCommand(cmd string, args ...string) (*exec.Cmd, error) {
	// TODO : log the command being run
	path, err := exec.LookPath(cmd)
	if err != nil {
		return nil, fmt.Errorf("executable not found: %w", err)
	}
	return exec.Command(path, args...), nil
}
