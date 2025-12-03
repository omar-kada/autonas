// Package storage defines all data store operations
package storage

import (
	"log/slog"
	"omar-kada/autonas/api"
)

// Storage is an abstraction of all database operations
type Storage interface {
	GetCurrentStacks() []string
	GetDeployments() ([]api.Deployment, error)
	SaveDeployment(deployment api.Deployment) (api.Deployment, error)
	UpdateStatus(deploymentID string, status api.DeploymentStatus) (api.Deployment, error)
	AddLogRecord(deploymentID string, record slog.Record) error
}
