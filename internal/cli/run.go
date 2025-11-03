// Package cli provides command-line interface functionalities for Autonas.
package cli

import (
	"fmt"
	"omar-kada/autonas/internal/config"
	"omar-kada/autonas/internal/exec"
	"omar-kada/autonas/internal/exec/git"
	"os"
	"reflect"

	"github.com/robfig/cron/v3"
)

var (
	generateConfigFromFiles = config.FromFiles
	syncCode                = git.SyncCode
	defaultDeployer         = exec.New()
)

// Runner defines the interface for running AutoNAS commands.
type Runner interface {
	RunCmd(configFiles []string, configRepo string) error
	RunPeriocically(cronPeriod string, configFiles []string, configRepo string)
}

// NewRunner creates a new Runner instance with default dependencies.
func NewRunner() Runner {
	return &runner{deployer: defaultDeployer}
}

type runner struct {
	deployer   exec.Deployer
	currentCfg config.Config
}

// RunCmd performs the main operations of fetching config, loading it, and deploying services.
func (r *runner) RunCmd(configFiles []string, configRepo string) error {

	cfg, err := generateConfigFromFiles(configFiles)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		return err
	}
	fmt.Printf("Final consolidated config: %+v\n", cfg)

	// TODO : add these to configuration
	configFolder := "."
	repoBranch := "main"

	err = syncCode(configRepo, repoBranch, configFolder)

	if err == git.NoErrAlreadyUpToDate {
		if reflect.DeepEqual(r.currentCfg, cfg) {
			fmt.Println("Configuration and repository are up to date. No changes detected.")
			return nil
		}
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting config repo: %v\n", err)
		return err
	}

	// Copy all files from ./services to SERVICES_PATH
	err = r.deployer.DeployServices(configFolder, r.currentCfg, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deploying services: %v\n", err)
		return err
	}
	r.currentCfg = cfg
	return nil
}

// RunPeriocically runs the RunCmd function periodically based on the given cron period string.
func (r *runner) RunPeriocically(cronPeriod string, configFiles []string, configRepo string) {
	c := cron.New()

	c.AddFunc(cronPeriod, func() {
		r.RunCmd(configFiles, configRepo)
	})

	c.Start()
	select {} // keep running
}
