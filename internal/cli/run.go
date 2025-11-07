// Package cli provides command-line interface functionalities for Autonas.
package cli

import (
	"fmt"
	"omar-kada/autonas/internal/config"
	"omar-kada/autonas/internal/containers"
	"omar-kada/autonas/internal/git"
	"omar-kada/autonas/internal/logger"
	"path/filepath"
	"reflect"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// Deployer abstracts service deployment operations
type Deployer interface {
	DeployServices(configFolder string, currentCfg, cfg config.Config) error
}

// New creates a new Runner instance with default dependencies.
func New(log logger.Logger) Runner {
	deployer := containers.NewDockerDeployer(log)
	return Runner{
		log:                      log,
		deployer:                 &deployer,
		_generateConfigFromFiles: config.FromFiles,
		_syncCode:                git.SyncCode,
	}
}

// Runner abstracts the implements of the run command
type Runner struct {
	log                      logger.Logger
	deployer                 Deployer
	_generateConfigFromFiles func(files []string) (config.Config, error)
	_syncCode                func(repoURL string, branch string, path string) error

	currentCfg config.Config
}

// RunCmd performs the main operations of fetching config, loading it, and deploying services.
func (r *Runner) RunCmd(configFiles []string, configRepo, configFolder string) error {
	// TODO : add these to configuration
	repoBranch := "main"
	syncErr := r._syncCode(configRepo, repoBranch, configFolder)

	if syncErr != nil && syncErr != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("error getting config repo:  %w", syncErr)

	}

	for i, file := range configFiles {
		if !filepath.IsAbs(file) {
			configFiles[i] = filepath.Join(configFolder, configFiles[i])
		}
	}
	cfg, err := r._generateConfigFromFiles(configFiles)

	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}
	r.log.Info("Final consolidated config", zap.Any("config", cfg))

	// check if the config changed from last run
	if syncErr == git.NoErrAlreadyUpToDate && reflect.DeepEqual(r.currentCfg, cfg) {
		r.log.Info("Configuration and repository are up to date. No changes detected.")
		return nil
	}

	// Copy all files from ./services to SERVICES_PATH
	err = r.deployer.DeployServices(configFolder, r.currentCfg, cfg)
	if err != nil {
		return fmt.Errorf("error deploying services: %w", err)
	}
	r.currentCfg = cfg
	return nil
}

// RunPeriodically runs the RunCmd function periodically based on the given cron period string.
func (r *Runner) RunPeriodically(cronPeriod string, configFiles []string, configRepo, configFolder string) {
	c := cron.New()

	c.AddFunc(cronPeriod, func() {
		err := r.RunCmd(configFiles, configRepo, configFolder)
		if err != nil {
			r.log.Errorf("error on run periodically: %w", err)
		}
	})

	c.Start()
	select {} // keep running
}
