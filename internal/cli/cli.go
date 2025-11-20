// Package cli handles command line operations
package cli

import (
	"github.com/spf13/cobra"
)

// NewRootCmd creates a new command with default dependencies
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "autonas",
		Short: "AutoNAS CLI",
	}
	rootCmd.AddCommand(newRunCommand())
	return rootCmd
}
