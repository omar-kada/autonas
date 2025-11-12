// Package model provides type definitions used with containers.
package model

import (
	"omar-kada/autonas/internal/config"
)

// Manager defines methods for managing containerized services.
type Manager interface {
	RemoveServices(services []string, servicesPath string) error
	DeployServices(cfg config.Config, servicesDir string) error
	// GetManagedContainers() (map[string][]Summary, error)
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
