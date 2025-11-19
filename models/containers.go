// Package models provides type definitions used with containers.
package models

// ContainerSummary is the domain view of a managed container.
type ContainerSummary struct {
	ID     string
	Names  []string
	Image  string
	State  string
	Status string
}
