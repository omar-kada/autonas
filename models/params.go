package models

import "os"

// DeploymentParams groups parameters related to the deployment process
type DeploymentParams struct {
	ConfigFile   string
	WorkingDir   string
	ServicesDir  string
	AddWritePerm bool
}

func (p DeploymentParams) GetAddWritePerm() os.FileMode {
	if p.AddWritePerm {
		return os.FileMode(0666)
	} else {
		return os.FileMode(0000)
	}
}

// ServerParams groups parameters related to the API server
type ServerParams struct {
	Port int
}
