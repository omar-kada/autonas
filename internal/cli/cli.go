// Package cli handles command line operations
package cli

import (
	"omar-kada/autonas/internal/storage"

	"github.com/spf13/cobra"
)

// NewRootCmd creates a new command with default dependencies
func NewRootCmd(store storage.Storage) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "autonas",
		Short: "AutoNAS CLI",
	}
	runCmd := newRunCmd(store)
	rootCmd.AddCommand(runCmd.ToCobraCommand())
	return rootCmd
}
