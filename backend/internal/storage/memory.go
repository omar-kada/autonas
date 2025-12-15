package storage

import (
	"fmt"
	"omar-kada/autonas/api"
	"omar-kada/autonas/models"
	"time"

	"github.com/docker/distribution/uuid"
	"github.com/elliotchance/orderedmap/v3"
)

// MemoryStorage uses memory to store data (to be used mainly for testing)
type MemoryStorage struct {
	deployments *orderedmap.OrderedMap[string, *api.Deployment]
}

// NewMemoryStorage instanciates a new memory storage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		deployments: orderedmap.NewOrderedMap[string, *api.Deployment](),
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
	for _, deployment := range s.deployments.AllFromBack() {
		deployments = append(deployments, *deployment)
	}

	return deployments, nil
}

// GetDeployment returns deployment by id
func (s *MemoryStorage) GetDeployment(id string) api.Deployment {
	return *s.deployments.GetOrDefault(id, nil)
}

func newID() string {
	return uuid.Generate().String()
}

// InitDeployment creates a new deployment and returns it
func (s *MemoryStorage) InitDeployment(title string, author string, diff string, files []api.FileDiff) (api.Deployment, error) {
	deployment := api.Deployment{
		Id:     newID(),
		Title:  title,
		Author: author,
		Time:   time.Now(),
		Status: "running",
		Diff:   diff,
		Files:  files,
		Events: []api.Event{},
	}
	s.deployments.Set(deployment.Id, &deployment)
	return deployment, nil
}

// UpdateStatus updates only the status of the deployment
func (s *MemoryStorage) UpdateStatus(deploymentID string, status api.DeploymentStatus) error {
	deployment, exists := s.deployments.Get(deploymentID)
	if !exists {
		return fmt.Errorf("deployment doesn't exist %s", deploymentID)
	}
	deployment.Status = status
	return nil
}

// StoreEvent saves the events to the corresponding deploymentID
func (s *MemoryStorage) StoreEvent(event models.Event) {
	s.deployments.GetOrDefault(event.ObjectID, nil).Events = append(s.deployments.GetOrDefault(event.ObjectID, nil).Events, api.Event{
		Level: api.EventLevel(event.Level.String()),
		Msg:   event.Msg,
		Time:  event.Time,
	})
}

// GetEvents retreives all events related to the deploymentID
func (s *MemoryStorage) GetEvents(objectID string) []api.Event {
	return s.deployments.GetOrDefault(objectID, &api.Deployment{}).Events
}
