// Package cli handles command line operations
package cli

import (
	"path/filepath"

	"omar-kada/autonas/internal/shell"
	"omar-kada/autonas/internal/storage"

	"github.com/spf13/cobra"
)

// NewRootCmd creates a new command with default dependencies
func NewRootCmd(executor shell.Executor) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "autonas",
		Short: "AutoNAS CLI",
	}
	rootCmd.AddCommand(NewRunCommand(executor, func(params RunParams) (storage.Storage, error) {
		return storage.NewGormStorage(
			filepath.Join(params.GetDBDir(), "autonas.db"),
			params.GetAddWritePerm(),
		)
	}))
	return rootCmd
}
