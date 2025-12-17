// Package models provides type definitions used with containers.
package models

import (
	"time"

	"github.com/moby/moby/api/types/container"
)

// ContainerSummary is the domain view of a managed container.
type ContainerSummary struct {
	ID        string
	Name      string
	Image     string
	State     container.ContainerState
	Health    container.HealthStatus
	StartedAt time.Time
}

// ContextKey is the type of keys used inside context
type ContextKey string
