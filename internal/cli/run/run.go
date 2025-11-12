// Package run provides the run command executor
package run

import (
	"fmt"
	"omar-kada/autonas/internal/cli/defaults"
	"omar-kada/autonas/internal/config"
	"omar-kada/autonas/internal/containers"
	"omar-kada/autonas/internal/git"
	"omar-kada/autonas/internal/logger"
	"path/filepath"
	"reflect"

	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const (
	_configDir   defaults.VarKey = "config-dir"
	_files       defaults.VarKey = "files"
	_branch      defaults.VarKey = "branch"
	_repo        defaults.VarKey = "repo"
	_workingDir  defaults.VarKey = "working-dir"
	_servicesDir defaults.VarKey = "services-dir"
	_cronPeriod  defaults.VarKey = "cron-period"
)

var varInfoMap = defaults.VariableInfoMap{
	_configDir:   {EnvKey: "AUTONAS_CONFIG_DIR", DefaultValue: nil},
	_files:       {EnvKey: "AUTONAS_CONFIG_FILES", DefaultValue: []string{"config.yaml"}},
	_repo:        {EnvKey: "AUTONAS_CONFIG_REPO", DefaultValue: nil},
	_branch:      {EnvKey: "AUTONAS_CONFIG_BRANCH", DefaultValue: "main"},
	_workingDir:  {EnvKey: "AUTONAS_WORKING_DIR", DefaultValue: "./config"},
	_servicesDir: {EnvKey: "AUTONAS_SERVICES_DIR", DefaultValue: "."},
	_cronPeriod:  {EnvKey: "AUTONAS_CRON_PERIOD", DefaultValue: nil},
}

var (
	getDefaultString  = defaults.GetDefaultStringFn(varInfoMap)
	envOrDefault      = defaults.EnvOrDefaultFn(varInfoMap)
	envOrDefaultSlice = defaults.EnvOrDefaultSliceFn(varInfoMap)
)

// Cmd abstracts the dependencies of the run command
type Cmd struct {
	Log             logger.Logger
	Deployer        containers.Deployer
	ConfigGenerator config.Generator
	Syncer          git.Syncer

	currentCfg config.Config
}

type runParams struct {
	ConfigFiles []string
	Repo        string
	Branch      string
	WorkingDir  string
	ServicesDir string
	CronPeriod  string
}

func getParamsWithDefaults(p runParams) runParams {
	return runParams{
		ConfigFiles: envOrDefaultSlice(p.ConfigFiles, _files),
		Repo:        envOrDefault(p.Repo, _repo),
		Branch:      envOrDefault(p.Branch, _branch),
		WorkingDir:  envOrDefault(p.WorkingDir, _workingDir),
		ServicesDir: envOrDefault(p.ServicesDir, _servicesDir),
		CronPeriod:  envOrDefault(p.CronPeriod, _cronPeriod),
	}
}

// New creates a new run.Cmd instance with default dependencies.
func New(log logger.Logger) Cmd {
	return Cmd{
		Log:             log,
		Deployer:        containers.NewDockerDeployer(log),
		ConfigGenerator: config.NewGenerator(),
		Syncer:          git.NewSyncer(),
	}
}

// ToCobraCommand transforms the run command to cobra.Command
func (r *Cmd) ToCobraCommand() *cobra.Command {
	params := runParams{}
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run with optional config files",
		Run: func(_ *cobra.Command, _ []string) {
			params := getParamsWithDefaults(params)

			r.RunOnce(params)
			if params.CronPeriod != "" {
				r.RunPeriodically(params)
			}
		},
	}

	runCmd.Flags().StringSliceVarP(&(params.ConfigFiles), string(_files), "f", nil,
		getDefaultString("YAML config files", _files))
	runCmd.Flags().StringVarP(&(params.Repo), string(_repo), "r", "",
		getDefaultString("repository URL to fetch config files & services", _repo))
	runCmd.Flags().StringVarP(&(params.Branch), string(_branch), "b", "",
		getDefaultString("branch to be used in the repo", _branch))
	runCmd.Flags().StringVarP(&(params.WorkingDir), string(_workingDir), "d", "",
		getDefaultString("directory where autonas data will be stored", _workingDir))
	runCmd.Flags().StringVarP(&(params.ServicesDir), string(_servicesDir), "s", "",
		getDefaultString("directory where services compose stacks will be stored", _servicesDir))
	runCmd.Flags().StringVarP(&(params.CronPeriod), string(_cronPeriod), "p", "",
		getDefaultString("cron period string", _cronPeriod))
	return runCmd
}

// RunOnce performs the main operations of fetching config, loading it, and deploying services.
func (r *Cmd) RunOnce(params runParams) error {
	syncErr := r.Syncer.Sync(params.Repo, params.Branch, params.WorkingDir)

	if syncErr != nil && syncErr != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("error getting config repo:  %w", syncErr)
	}

	for i, file := range params.ConfigFiles {
		if !filepath.IsAbs(file) {
			params.ConfigFiles[i] = filepath.Join(params.WorkingDir, params.ConfigFiles[i])
		}
	}
	cfg, err := r.ConfigGenerator.FromFiles(params.ConfigFiles)
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}
	r.Log.Debug("Final consolidated config", zap.Any("config", cfg))

	// check if the config changed from last run
	if syncErr == git.NoErrAlreadyUpToDate && reflect.DeepEqual(r.currentCfg, cfg) {
		r.Log.Info("Configuration and repository are up to date. No changes detected.")
		return nil
	}

	// Copy all files from ./services to SERVICES_PATH
	err = r.Deployer.DeployServices(params.WorkingDir, params.ServicesDir, r.currentCfg, cfg)
	if err != nil {
		return fmt.Errorf("error deploying services: %w", err)
	}
	r.currentCfg = cfg
	return nil
}

// RunPeriodically runs the RunCmd function periodically based on the given cron period string.
func (r *Cmd) RunPeriodically(params runParams) {
	c := cron.New()

	c.AddFunc(params.CronPeriod, func() {
		err := r.RunOnce(params)
		if err != nil {
			r.Log.Errorf("error on run periodically: %w", err)
		}
	})

	c.Start()
	select {} // keep running
}
