package models

import (
	"os"
	"path/filepath"
)

// DeploymentParams groups parameters related to the deployment process
type DeploymentParams struct {
	WorkingDir   string
	ServicesDir  string
	AddWritePerm bool
}

// GetAddWritePerm returns the permissions to add based on the AddWritePerm boolean
func (p DeploymentParams) GetAddWritePerm() os.FileMode {
	if p.AddWritePerm {
		return os.FileMode(0666)
	}
	return os.FileMode(0000)
}

func (p DeploymentParams) GetRepoDir() string {
	return filepath.Join(p.WorkingDir, "repo")
}

// ServerParams groups parameters related to the API server
type ServerParams struct {
	Port int
}
