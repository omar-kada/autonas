// Package models provides type definitions used with containers.
package models

import "github.com/moby/moby/api/types/container"

// ContainerSummary is the domain view of a managed container.
type ContainerSummary struct {
	ID     string
	Name   string
	Image  string
	State  container.ContainerState
	Health container.HealthStatus
}
