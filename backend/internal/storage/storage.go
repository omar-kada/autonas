// Package storage defines all data store operations
package storage

import (
	"omar-kada/autonas/models"
)

// Storage is an abstraction of all database operations
type Storage interface {
	DeploymentStorage
	EventStorage
}

// DeploymentStorage is an abstraction of all deployment database operations
type DeploymentStorage interface {
	GetDeployments(c Cursor[uint64]) ([]models.Deployment, error)
	GetDeployment(id uint64) (models.Deployment, error)
	InitDeployment(title string, author string, diff string, files []models.FileDiff) (models.Deployment, error)
	EndDeployment(deploymentID uint64, status models.DeploymentStatus) error
	GetLastDeployment() (models.Deployment, error)
}

// EventStorage is an abstraction of all event database operations
type EventStorage interface {
	StoreEvent(event models.Event) error
	GetEvents(objectID uint64) ([]models.Event, error)
}
