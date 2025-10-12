package main

import (
	"fmt"
	"os"
	"omar-kada/autonas/internal/cli"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "autonas",
	Short: "AutoNAS CLI",
}

var (
	configFiles []string
	configRepo string
	runCmd      *cobra.Command
)

func init() {
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run with optional config files",
		Run: func(cmd *cobra.Command, args []string) {
			cli.RunCmd(configFiles, configRepo)
		},
	}
	runCmd.Flags().StringSliceVarP(&configFiles, "config", "c", []string{"config.default.yaml", "config.yaml"}, "YAML config files (default: config.yaml)")
	runCmd.Flags().StringVarP(&configRepo, "repo", "r","", "repository URL to fetch config files & services")
}
func main() {
	// Add subcommands
	rootCmd.AddCommand(runCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
