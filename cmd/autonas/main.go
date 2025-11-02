// Package main is the entry point for AutoNAS.
package main

import (
	"fmt"
	"omar-kada/autonas/internal/cli"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "autonas",
	Short: "AutoNAS CLI",
}

var (
	configFiles []string
	configRepo  string
	cronPeriod  string
	runCmd      *cobra.Command
)

func init() {
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run with optional config files",
		Run: func(_ *cobra.Command, _ []string) {
			runner := cli.NewRunner()
			runner.RunCmd(configFiles, configRepo)
			if cronPeriod != "" {
				runner.RunPeriocically(cronPeriod, configFiles, configRepo)
			}
		},
	}
	runCmd.Flags().StringSliceVarP(&configFiles, "config", "c", []string{"config.default.yaml", "config.yaml"}, "YAML config files (default: config.yaml)")
	runCmd.Flags().StringVarP(&configRepo, "repo", "r", "", "repository URL to fetch config files & services")
	runCmd.Flags().StringVarP(&cronPeriod, "period", "p", "", "cron period string")
}
func main() {
	// Add subcommands
	rootCmd.AddCommand(runCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
