package cli

import (
	"fmt"
	"log/slog"
	"omar-kada/autonas/internal/api"
	"omar-kada/autonas/internal/cli/defaults"
	"omar-kada/autonas/internal/config"
	"omar-kada/autonas/internal/docker"
	"omar-kada/autonas/internal/files"
	"omar-kada/autonas/internal/git"
	"omar-kada/autonas/internal/process"
	"omar-kada/autonas/internal/storage"
	"omar-kada/autonas/models"

	"github.com/moby/moby/client"
	"github.com/spf13/cobra"
)

const (
	_file         defaults.VarKey = "file"
	_workingDir   defaults.VarKey = "working-dir"
	_servicesDir  defaults.VarKey = "services-dir"
	_addWritePerm defaults.VarKey = "add-write-perm"
	_port         defaults.VarKey = "port"
)

var varInfoMap = defaults.VariableInfoMap{
	_file:         {EnvKey: "AUTONAS_CONFIG_FILE", DefaultValue: "/config/config.yaml"},
	_workingDir:   {EnvKey: "AUTONAS_WORKING_DIR", DefaultValue: "./config"},
	_servicesDir:  {EnvKey: "AUTONAS_SERVICES_DIR", DefaultValue: "."},
	_addWritePerm: {DefaultValue: false},
	_port:         {EnvKey: "AUTONAS_PORT", DefaultValue: 5005},
}

// RunParams contain parameters of the run command
type RunParams struct {
	models.DeploymentParams
	models.ServerParams
}

func getParamsWithDefaults(p RunParams) RunParams {
	return RunParams{
		DeploymentParams: models.DeploymentParams{
			ConfigFile:   varInfoMap.EnvOrDefault(p.ConfigFile, _file),
			WorkingDir:   varInfoMap.EnvOrDefault(p.WorkingDir, _workingDir),
			ServicesDir:  varInfoMap.EnvOrDefault(p.ServicesDir, _servicesDir),
			AddWritePerm: p.AddWritePerm,
		},
		ServerParams: models.ServerParams{
			Port: varInfoMap.EnvOrDefaultInt(p.Port, _port),
		},
	}
}

// newRunCommand creates a new run
func newRunCommand() *cobra.Command {
	params := RunParams{}
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run with optional config files",
		Run: func(_ *cobra.Command, _ []string) {
			doRun(getParamsWithDefaults(params))
		},
	}

	runCmd.Flags().StringVarP(&(params.ConfigFile), string(_file), "f", "",
		varInfoMap.GetDefaultString("YAML config file", _file))
	runCmd.Flags().StringVarP(&(params.WorkingDir), string(_workingDir), "d", "",
		varInfoMap.GetDefaultString("directory where autonas data will be stored", _workingDir))
	runCmd.Flags().StringVarP(&(params.ServicesDir), string(_servicesDir), "s", "",
		varInfoMap.GetDefaultString("directory where services compose stacks will be stored", _servicesDir))
	runCmd.Flags().BoolVar(&(params.AddWritePerm), string(_addWritePerm), false,
		varInfoMap.GetDefaultString("when true, the tool adds write permission to config files", _addWritePerm))
	runCmd.Flags().IntVarP(&(params.Port), string(_port), "p", 5005,
		varInfoMap.GetDefaultString("port that will be used for exposing the API", _port))
	return runCmd
}

func doRun(params RunParams) error {
	store := storage.NewMemoryStorage()
	dockerClient, err := client.New(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("failed to create docker client: %w", err)
	}

	manager := process.NewManager(
		params.DeploymentParams,
		docker.NewDeployer(),
		docker.NewInspector(dockerClient),
		files.NewCopier(),
		git.NewFetcher(),
		config.NewGenerator())

	if err := manager.SyncDeployment(); err != nil {
		slog.Error("error while deploying services", "error", err)
		return err
	}
	server := api.NewServer(store, manager)
	return server.ListenAndServe(params.Port)
}
