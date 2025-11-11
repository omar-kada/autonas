// Package cli handles command line operations
package cli

import (
	"omar-kada/autonas/internal/cli/run"
	"omar-kada/autonas/internal/logger"

	"github.com/spf13/cobra"
)

// NewDefaultRunCmd creates a new "run" command with default dependencies
func NewDefaultRunCmd(log logger.Logger) *cobra.Command {
	runCmd := run.New(log)
	return runCmd.ToCobraCommand()
}
