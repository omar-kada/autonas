// Package storage defines all data store operations
package storage

import (
	"omar-kada/autonas/api"
	"omar-kada/autonas/models"
)

// Storage is an abstraction of all database operations
type Storage interface {
	DeploymentStorage
	EventStorage
}

// DeploymentStorage is an abstraction of all deployment database operations
type DeploymentStorage interface {
	GetCurrentStacks() []string
	GetDeployments() ([]api.Deployment, error)
	GetDeployment(id string) api.Deployment
	InitDeployment(title string, author string, diff string, files []api.FileDiff) (api.Deployment, error)
	UpdateStatus(deploymentID string, status api.DeploymentStatus) error
}

// EventStorage is an abstraction of all event database operations
type EventStorage interface {
	StoreEvent(event models.Event)
	GetEvents(objectID string) []api.Event
}
