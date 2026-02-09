package mappers

import (
	"omar-kada/autonas/api"
	"omar-kada/autonas/models"
)

// StatusMapper maps models.ContainerSummary to api.ContainerStatus
type StatusMapper struct{}

// Map converts a models.ContainerSummary to an api.ContainerStatus
func (StatusMapper) Map(container models.ContainerSummary) api.ContainerStatus {
	return api.ContainerStatus{
		ContainerId: container.ID,
		State:       api.ContainerStatusState(container.State),
		Name:        container.Name,
		Health:      api.ContainerStatusHealth(container.Health),
		StartedAt:   container.StartedAt,
	}
}
