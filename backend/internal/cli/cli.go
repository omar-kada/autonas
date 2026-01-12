// Package cli handles command line operations
package cli

import (
	"omar-kada/autonas/internal/shell"

	"github.com/spf13/cobra"
)

// NewRootCmd creates a new command with default dependencies
func NewRootCmd(executor shell.Executor) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "autonas",
		Short: "AutoNAS CLI",
	}
	rootCmd.AddCommand(NewRunCommand(executor))
	return rootCmd
}
