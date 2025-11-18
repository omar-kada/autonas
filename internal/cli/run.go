package cli

import (
	"fmt"
	"log/slog"
	"omar-kada/autonas/internal/api"
	"omar-kada/autonas/internal/cli/defaults"
	"omar-kada/autonas/internal/config"
	"omar-kada/autonas/internal/containers"
	"omar-kada/autonas/internal/git"
	"omar-kada/autonas/internal/storage"
	"reflect"

	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

const (
	_configDir    defaults.VarKey = "config-dir"
	_files        defaults.VarKey = "files"
	_branch       defaults.VarKey = "branch"
	_repo         defaults.VarKey = "repo"
	_workingDir   defaults.VarKey = "working-dir"
	_servicesDir  defaults.VarKey = "services-dir"
	_cronPeriod   defaults.VarKey = "cron-period"
	_addWritePerm defaults.VarKey = "add-write-perm"
	_port         defaults.VarKey = "port"
)

var varInfoMap = defaults.VariableInfoMap{
	_configDir:    {EnvKey: "AUTONAS_CONFIG_DIR", DefaultValue: nil},
	_files:        {EnvKey: "AUTONAS_CONFIG_FILES", DefaultValue: []string{"config.yaml"}},
	_repo:         {EnvKey: "AUTONAS_CONFIG_REPO", DefaultValue: nil},
	_branch:       {EnvKey: "AUTONAS_CONFIG_BRANCH", DefaultValue: "main"},
	_workingDir:   {EnvKey: "AUTONAS_WORKING_DIR", DefaultValue: "./config"},
	_servicesDir:  {EnvKey: "AUTONAS_SERVICES_DIR", DefaultValue: "."},
	_cronPeriod:   {EnvKey: "AUTONAS_CRON_PERIOD", DefaultValue: nil},
	_addWritePerm: {DefaultValue: false},
	_port:         {DefaultValue: 8080},
}

// Cmd abstracts the dependencies of the run command
type Cmd struct {
	deployer        containers.Deployer
	configGenerator config.Generator
	syncer          git.Syncer
	store           storage.Storage
	server          api.Server

	currentCfg config.Config
}

type runParams struct {
	ConfigFiles  []string
	Repo         string
	Branch       string
	WorkingDir   string
	ServicesDir  string
	CronPeriod   string
	AddWritePerm bool
	Port         int
}

func getParamsWithDefaults(p runParams) runParams {
	return runParams{
		ConfigFiles:  varInfoMap.EnvOrDefaultSlice(p.ConfigFiles, _files),
		Repo:         varInfoMap.EnvOrDefault(p.Repo, _repo),
		Branch:       varInfoMap.EnvOrDefault(p.Branch, _branch),
		WorkingDir:   varInfoMap.EnvOrDefault(p.WorkingDir, _workingDir),
		ServicesDir:  varInfoMap.EnvOrDefault(p.ServicesDir, _servicesDir),
		CronPeriod:   varInfoMap.EnvOrDefault(p.CronPeriod, _cronPeriod),
		AddWritePerm: p.AddWritePerm,
		Port:         varInfoMap.EnvOrDefaultInt(p.Port, _port),
	}
}

func newRunCmd(store storage.Storage) Cmd {
	return Cmd{
		deployer:        containers.NewDockerDeployer(),
		configGenerator: config.NewGenerator(),
		syncer:          git.NewSyncer(),
		store:           store,
		server:          api.NewServer(store),
	}
}

// ToCobraCommand transforms the run command to cobra.Command
func (r *Cmd) ToCobraCommand() *cobra.Command {
	params := runParams{}
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run with optional config files",
		Run: func(_ *cobra.Command, _ []string) {
			r.DoRun(getParamsWithDefaults(params))
		},
	}

	runCmd.Flags().StringSliceVarP(&(params.ConfigFiles), string(_files), "f", nil,
		varInfoMap.GetDefaultString("YAML config files", _files))
	runCmd.Flags().StringVarP(&(params.Repo), string(_repo), "r", "",
		varInfoMap.GetDefaultString("repository URL to fetch config files & services", _repo))
	runCmd.Flags().StringVarP(&(params.Branch), string(_branch), "b", "",
		varInfoMap.GetDefaultString("branch to be used in the repo", _branch))
	runCmd.Flags().StringVarP(&(params.WorkingDir), string(_workingDir), "d", "",
		varInfoMap.GetDefaultString("directory where autonas data will be stored", _workingDir))
	runCmd.Flags().StringVarP(&(params.ServicesDir), string(_servicesDir), "s", "",
		varInfoMap.GetDefaultString("directory where services compose stacks will be stored", _servicesDir))
	runCmd.Flags().StringVarP(&(params.CronPeriod), string(_cronPeriod), "p", "",
		varInfoMap.GetDefaultString("cron period string", _cronPeriod))
	runCmd.Flags().BoolVar(&(params.AddWritePerm), string(_addWritePerm), false,
		varInfoMap.GetDefaultString("when true, the tool adds write permission to config files", _addWritePerm))
	return runCmd
}

// DoRun executes the run command based on the input params
func (r *Cmd) DoRun(params runParams) {
	if params.AddWritePerm {
		r.deployer = r.deployer.WithPermission(0666)
	}
	r.RunOnce(params)
	if params.CronPeriod != "" {
		go r.RunPeriodically(params)
	}

	go r.server.ListenAndServe(params.Port)
	select {}
}

// RunOnce performs the main operations of fetching config, loading it, and deploying services.
func (r *Cmd) RunOnce(params runParams) error {
	syncErr := r.syncer.Sync(params.Repo, params.Branch, params.WorkingDir)

	if syncErr != nil && syncErr != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("error getting config repo:  %w", syncErr)
	}

	cfg, err := r.configGenerator.FromFiles(params.ConfigFiles)
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}
	slog.Debug("Final consolidated config", "config", cfg)

	// check if the config changed from last run
	if syncErr == git.NoErrAlreadyUpToDate && reflect.DeepEqual(r.currentCfg, cfg) {
		slog.Info("Configuration and repository are up to date. No changes detected.")
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
func (r *Cmd) RunPeriodically(params runParams) {
	c := cron.New()

	c.AddFunc(params.CronPeriod, func() {
		err := r.RunOnce(params)
		if err != nil {
			slog.Error("error on run periodically", "error", err)
		}
	})

	c.Start()
	select {} // keep running
}
