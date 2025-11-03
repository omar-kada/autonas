// Package containers provides type definitions used with containers.
package containers

import (
	"omar-kada/autonas/internal/config"
)

// Handler defines methods for managing containerized services.
type Handler interface {
	RemoveServices(services []string, servicesPath string) error
	DeployServices(cfg config.Config) error
	GetManagedContainers() (map[string][]Summary, error)
}

// New creates a new instance of the Handler.
func New() Handler {
	return newDockerHandler()
}

// Summary is the domain view of a managed container.
// Keep only fields callers need.
type Summary struct {
	ID     string
	Names  []string
	Image  string
	State  string
	Status string
}
