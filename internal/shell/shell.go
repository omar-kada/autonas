// Package shell provides utilities to interact with shell.
package shell

import (
	"os"
	"os/exec"
	"runtime"
)

// RunShellCommand runs a shell command and returns error if any
func RunShellCommand(cmdStr string) error {
	var c *exec.Cmd
	if runtime.GOOS == "windows" {
		c = execCommand("cmd", "/C", cmdStr)
	} else {
		shell := os.Getenv("SHELL")
		if shell == "" {
			shell = "bash"
		}
		c = execCommand(shell, "-c", cmdStr)
	}
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

// execCommand is a wrapper for exec.Command for testability
var execCommand = defaultExecCommand

func defaultExecCommand(name string, arg ...string) *exec.Cmd {
	return exec.Command(name, arg...)
}
