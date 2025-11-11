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

	// Add subcommands
	rootCmd.AddCommand(cli.NewDefaultRunCmd(log))

	if err := rootCmd.Execute(); err != nil {
		log.Errorf("error on the root command : %w", err)
		os.Exit(1)
	}
}
