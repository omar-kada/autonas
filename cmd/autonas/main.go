// Package main is the entry point for AutoNAS.
package main

import (
	"omar-kada/autonas/internal/cli"
	"omar-kada/autonas/internal/logger"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "autonas",
	Short: "AutoNAS CLI",
}

var (
	configFiles []string
	configRepo  string
	directory   string
	cronPeriod  string
	log         logger.Logger
	runCmd      *cobra.Command
)

func main() {

	env := strings.ToUpper(os.Getenv("ENV"))
	log = logger.New(env == "DEV")
	defer log.Sync()

	runner := cli.New(log)
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run with optional config files",
		Run: func(_ *cobra.Command, _ []string) {
			runner.RunCmd(configFiles, configRepo, directory)
			if cronPeriod != "" {
				runner.RunPeriodically(cronPeriod, configFiles, configRepo, directory)
			}
		},
	}
	runCmd.Flags().StringSliceVarP(&configFiles, "config", "c", []string{"config.yaml"}, "YAML config files (default: config.yaml)")
	runCmd.Flags().StringVarP(&configRepo, "repo", "r", "", "repository URL to fetch config files & services")
	runCmd.Flags().StringVarP(&directory, "directory", "d", ".", "repo clone directory (default : '.'")
	runCmd.Flags().StringVarP(&cronPeriod, "period", "p", "", "cron period string")

	// Add subcommands
	rootCmd.AddCommand(runCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Error("error on the root command : %w", err)
		os.Exit(1)
	}
}
