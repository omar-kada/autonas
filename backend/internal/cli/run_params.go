package cli

import (
	"omar-kada/autonas/internal/cli/defaults"
	"omar-kada/autonas/models"
)

const (
	_file         defaults.VarKey = "file"
	_workingDir   defaults.VarKey = "working-dir"
	_servicesDir  defaults.VarKey = "services-dir"
	_addWritePerm defaults.VarKey = "add-write-perm"
	_port         defaults.VarKey = "port"
)

var varInfoMap = defaults.VariableInfoMap{
	_file:         {EnvKey: "AUTONAS_CONFIG_FILE", DefaultValue: "/data/config.yaml"},
	_workingDir:   {EnvKey: "AUTONAS_WORKING_DIR", DefaultValue: "./config"},
	_servicesDir:  {EnvKey: "AUTONAS_SERVICES_DIR", DefaultValue: "."},
	_addWritePerm: {EnvKey: "AUTONAS_ADD_WRITE_PERM", DefaultValue: "false"},
	_port:         {EnvKey: "AUTONAS_PORT", DefaultValue: 0},
}

// RunParams contain parameters of the run command
type RunParams struct {
	models.DeploymentParams
	models.ServerParams
	ConfigFile string
}

func getParamsWithDefaults(p RunParams) RunParams {
	return RunParams{
		ConfigFile: varInfoMap.EnvOrDefault(p.ConfigFile, _file),
		DeploymentParams: models.DeploymentParams{
			WorkingDir:   varInfoMap.EnvOrDefault(p.WorkingDir, _workingDir),
			ServicesDir:  varInfoMap.EnvOrDefault(p.ServicesDir, _servicesDir),
			AddWritePerm: varInfoMap.EnvOrDefault(p.AddWritePerm, _addWritePerm),
		},
		ServerParams: models.ServerParams{
			Port: varInfoMap.EnvOrDefaultInt(p.Port, _port),
		},
	}
}
