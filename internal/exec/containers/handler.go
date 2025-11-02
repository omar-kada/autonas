// Package containers provides utilities to interact with Container managers.
package containers

import (
	"omar-kada/autonas/internal/config"

	"github.com/moby/moby/api/types/container"
)

// Handler defines methods for managing containerized services.
type Handler interface {
	RemoveServices(services []string, servicesPath string) error
	DeployServices(cfg config.Config) error
	GetManagedContainers() (map[string][]container.Summary, error)
}

// NewHandler creates a new instance of the Containers Handler.
func NewHandler() Handler {
	return newDockerHandler()
}
