// Package storage defines all data store operations
package storage

import (
	"omar-kada/autonas/modelsdb"
)

// Storage is an abstraction of all database operations
type Storage interface {
	DeploymentStorage
	EventStorage
}

// DeploymentStorage is an abstraction of all deployment database operations
type DeploymentStorage interface {
	GetDeployments() ([]*modelsdb.Deployment, error)
	GetDeployment(id uint64) (*modelsdb.Deployment, error)
	InitDeployment(title string, author string, diff string, files []*modelsdb.FileDiff) (*modelsdb.Deployment, error)
	EndDeployment(deploymentID uint64, status modelsdb.DeploymentStatus) error
}

// EventStorage is an abstraction of all event database operations
type EventStorage interface {
	StoreEvent(event modelsdb.Event) error
	GetEvents(objectID uint64) ([]*modelsdb.Event, error)
}
