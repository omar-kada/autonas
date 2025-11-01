// Package cli provides command-line interface functionalities for Autonas.
package cli

import (
	"fmt"
	"omar-kada/autonas/internal/config"
	"omar-kada/autonas/internal/exec"
	"omar-kada/autonas/internal/exec/git"
	"os"

	"github.com/robfig/cron/v3"
)

var (
	generateConfigFromFiles = config.FromFiles
	syncCode                = git.SyncCode
	deployServices          = exec.DeployServices
)

var currentCfg config.Config

// RunCmd performs the main operations of fetching config, loading it, and deploying services.
func RunCmd(configFiles []string, configRepo string) error {

	// TODO : add these to configuration
	configFolder := "."
	repoBranch := "main"

	err := syncCode(configRepo, repoBranch, configFolder)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting config repo: %v\n", err)
		return err
	}

	cfg, err := generateConfigFromFiles(configFiles)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		return err
	}
	fmt.Printf("Final consolidated config: %+v\n", cfg)

	// Copy all files from ./services to SERVICES_PATH
	err = deployServices(configFolder, currentCfg, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deploying services: %v\n", err)
		return err
	}
	currentCfg = cfg
	return nil
}

// RunPeriocically runs the RunCmd function periodically based on the given cron period string.
func RunPeriocically(cronPeriod string, configFiles []string, configRepo string) {
	c := cron.New()

	c.AddFunc(cronPeriod, func() {
		RunCmd(configFiles, configRepo)
	})

	c.Start()
	select {} // keep running
}
