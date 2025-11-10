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
	branch      string
	workingDir  string
	servicesDir string
	cronPeriod  string
)

/*

tool config (cmd args, or ENV variables):
	- files : list of config files names [default : config.yaml]
	- repo : config Repo [Required]
	- branch : repo branch [default: main]
	- services-directory : target services folder [Required]
	- working-directory : directory where temp files will be created [default "."]
	- cron : CRON schedule [Optional]

deployment config files (everything reltaed to deployed services):
	- services-directory : target services folder
	- data-directory : target data folder
	- global variables available to all services
	- enabledServices
	- different services specific variables

*/

func main() {

	env := strings.ToUpper(os.Getenv("ENV"))
	log := logger.New(env == "DEV")
	defer log.Sync()

	runner := cli.New(log)
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run with optional config files",
		Run: func(_ *cobra.Command, _ []string) {
			params := cli.RunParams{
				ConfigFiles: configFiles,
				Repo:        configRepo,
				Branch:      branch,
				WorkingDir:  workingDir,
				ServicesDir: servicesDir,
			}
			runner.RunCmd(params)
			if cronPeriod != "" {
				runner.RunPeriodically(cronPeriod, params)
			}
		},
	}
	runCmd.Flags().StringSliceVarP(
		&configFiles, "files", "f",
		envOrDefaultSlice("AUTONAS_CONFIG_FILES", []string{"config.yaml"}),
		"YAML config files (default: config.yaml)",
	)
	runCmd.Flags().StringVarP(
		&configRepo, "repo", "r",
		envOrDefault("AUTONAS_CONFIG_REPO", ""),
		"repository URL to fetch config files & services",
	)
	runCmd.Flags().StringVarP(
		&branch, "branch", "b",
		envOrDefault("AUTONAS_CONFIG_BRANCH", ""),
		"branch to be used in the repo",
	)
	runCmd.Flags().StringVarP(
		&workingDir, "working-directory", "d",
		envOrDefault("AUTONAS_WORKING_DIRECTORY", "."),
		"directory where autonas data will be stored (default : '.'",
	)
	runCmd.Flags().StringVarP(
		&servicesDir, "services-directory", "s",
		envOrDefault("AUTONAS_SERVICES_DIRECTORY", "."),
		"directory where services compose stacks will be stored (default : '.'",
	)
	runCmd.Flags().StringVarP(
		&cronPeriod, "cron-period", "c",
		envOrDefault("AUTONAS_CRON_PERIOD", ""),
		"cron period string",
	)

	// Add subcommands
	rootCmd.AddCommand(runCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Errorf("error on the root command : %w", err)
		os.Exit(1)
	}
}

func envOrDefault(envVar string, defaultValue string) string {
	value := os.Getenv(envVar)
	if value == "" {
		value = defaultValue
	}
	return value
}

func envOrDefaultSlice(envVar string, defaultValue []string) []string {
	value := os.Getenv(envVar)
	if value == "" {
		return defaultValue
	}
	return strings.Split(value, ",")
}
