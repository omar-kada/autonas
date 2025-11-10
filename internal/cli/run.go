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
	DeployServices(configFolder, servicesDir string, currentCfg, cfg config.Config) error
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

// RunParams contains all input parameters of the run command
type RunParams struct {
	ConfigFiles []string
	Repo        string
	Branch      string
	WorkingDir  string
	ServicesDir string
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
func (r *Runner) RunCmd(params RunParams) error {

	syncErr := r._syncCode(params.Repo, params.Branch, params.WorkingDir)

	if syncErr != nil && syncErr != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("error getting config repo:  %w", syncErr)

	}

	for i, file := range params.ConfigFiles {
		if !filepath.IsAbs(file) {
			params.ConfigFiles[i] = filepath.Join(params.WorkingDir, params.ConfigFiles[i])
		}
	}
	cfg, err := r._generateConfigFromFiles(params.ConfigFiles)

	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}
	r.log.Debug("Final consolidated config", zap.Any("config", cfg))

	// check if the config changed from last run
	if syncErr == git.NoErrAlreadyUpToDate && reflect.DeepEqual(r.currentCfg, cfg) {
		r.log.Info("Configuration and repository are up to date. No changes detected.")
		return nil
	}

	// Copy all files from ./services to SERVICES_PATH
	err = r.deployer.DeployServices(params.WorkingDir, params.ServicesDir, r.currentCfg, cfg)
	if err != nil {
		return fmt.Errorf("error deploying services: %w", err)
	}
	r.currentCfg = cfg
	return nil
}

// RunPeriodically runs the RunCmd function periodically based on the given cron period string.
func (r *Runner) RunPeriodically(cronPeriod string, params RunParams) {
	c := cron.New()

	c.AddFunc(cronPeriod, func() {
		err := r.RunCmd(params)
		if err != nil {
			r.log.Errorf("error on run periodically: %w", err)
		}
	})

	c.Start()
	select {} // keep running
}
