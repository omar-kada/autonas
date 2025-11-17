// Package main is the entry point for AutoNAS.
package main

import (
	"omar-kada/autonas/internal/cli"
	"omar-kada/autonas/internal/logger"
	"omar-kada/autonas/internal/storage"
	"os"
	"strings"
)

func main() {
	retcode := 0
	defer func() { os.Exit(retcode) }()

	log := logger.New(strings.ToUpper(os.Getenv("ENV")) == "DEV")
	defer log.Sync()

	store := storage.NewMemoryStorage()

	// Add subcommands
	rootCmd := cli.NewRootCmd(store, log)
	if err := rootCmd.Execute(); err != nil {
		log.Errorf("error on the root command : %w", err)
		retcode = 1 // it exits with code 1
	}
}
