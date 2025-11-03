// Package shell provides utilities to interact with shell.
package shell

import (
	"fmt"
	"os"
	"os/exec"
)

// RunCommand runs a shell command and returns error if any
func RunCommand(cmdAndArgs ...string) error {
	path, err := exec.LookPath(cmdAndArgs[0])
	if err != nil {
		return fmt.Errorf("executable not found: %w", err)
	}
	c, err := execCommand(path, cmdAndArgs[1:]...)
	if err != nil {
		return err
	}
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	fmt.Printf("Running command: %s\n", c)
	return c.Run()
}

// execCommand is a wrapper for exec.Command for testability
var execCommand = defaultExecCommand

func defaultExecCommand(cmd string, args ...string) (*exec.Cmd, error) {
	return exec.Command(cmd, args...), nil
}
