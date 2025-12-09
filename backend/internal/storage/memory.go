package storage

import (
	"fmt"
	"omar-kada/autonas/api"
	"omar-kada/autonas/models"
	"time"

	"github.com/docker/distribution/uuid"
)

// MemoryStorage uses memory to store data (to be used mainly for testing)
type MemoryStorage struct {
	deployments map[string]*api.Deployment
}

// NewMemoryStorage instanciates a new memory storage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		deployments: make(map[string]*api.Deployment),
	}
}

// GetCurrentStacks returns currently stored stacks
func (*MemoryStorage) GetCurrentStacks() []string {
	return []string{"test"}
}

// GetDeployments returns stored deployments
func (s *MemoryStorage) GetDeployments() ([]api.Deployment, error) {
	// transform map to slice
	var deployments []api.Deployment
	for _, deployment := range s.deployments {
		deployments = append(deployments, *deployment)
	}

	return deployments, nil
}

func newID() string {
	return uuid.Generate().String()
}

// InitDeployment creates a new deployment and returns it
func (s *MemoryStorage) InitDeployment(title string, diff string) (api.Deployment, error) {
	deployment := api.Deployment{
		Id:     newID(),
		Title:  title,
		Time:   time.Now(),
		Status: "running",
		Diff:   diff,
		Events: []api.Event{},
	}
	s.deployments[deployment.Id] = &deployment
	return deployment, nil
}

// UpdateStatus updates only the status of the deployment
func (s *MemoryStorage) UpdateStatus(deploymentID string, status api.DeploymentStatus) error {
	deployment, exists := s.deployments[deploymentID]
	if !exists {
		return fmt.Errorf("deployment doesn't exist %s", deploymentID)
	}
	deployment.Status = status
	return nil
}

// StoreEvent saves the events to the corresponding deploymentID
func (s *MemoryStorage) StoreEvent(event models.Event) {
	s.deployments[event.ObjectID].Events = append(s.deployments[event.ObjectID].Events, api.Event{
		Level: api.EventLevel(event.Level.String()),
		Msg:   event.Msg,
		Time:  event.Time,
	})
}

// GetEvents retreives all events related to the deploymentID
func (s *MemoryStorage) GetEvents(objectID string) []api.Event {
	return s.deployments[objectID].Events
}
