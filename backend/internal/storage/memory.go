package storage

import (
	"fmt"
	"omar-kada/autonas/modelsdb"
	"time"

	"github.com/elliotchance/orderedmap/v3"
)

// MemoryStorage uses memory to store data (to be used mainly for testing)
type MemoryStorage struct {
	deployments *orderedmap.OrderedMap[uint64, *modelsdb.Deployment]
}

// NewMemoryStorage instanciates a new memory storage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		deployments: orderedmap.NewOrderedMap[uint64, *modelsdb.Deployment](),
	}
}

// GetDeployments returns stored deployments
func (s *MemoryStorage) GetDeployments() ([]*modelsdb.Deployment, error) {
	// transform map to slice
	var deployments []*modelsdb.Deployment
	for _, deployment := range s.deployments.AllFromBack() {
		deployments = append(deployments, deployment)
	}

	return deployments, nil
}

// GetDeployment returns deployment by id
func (s *MemoryStorage) GetDeployment(id uint64) (*modelsdb.Deployment, error) {
	deployment, exists := s.deployments.Get(id)
	if !exists {
		return nil, fmt.Errorf("deployment doesn't exist %d", id)
	}
	return deployment, nil
}

func newID() uint64 {
	id := uint64(time.Now().UnixNano())
	return id
}

// InitDeployment creates a new deployment and returns it
func (s *MemoryStorage) InitDeployment(title string, author string, diff string, files []*modelsdb.FileDiff) (*modelsdb.Deployment, error) {
	deployment := modelsdb.Deployment{
		ID:     newID(),
		Title:  title,
		Author: author,
		Time:   time.Now(),
		Status: "running",
		Diff:   diff,
		Files:  files,
		Events: []*modelsdb.Event{},
	}
	s.deployments.Set(deployment.ID, &deployment)
	return &deployment, nil
}

// EndDeployment updates only the status of the deployment
func (s *MemoryStorage) EndDeployment(deploymentID uint64, status modelsdb.DeploymentStatus) error {
	deployment, exists := s.deployments.Get(deploymentID)
	if !exists {
		return fmt.Errorf("deployment doesn't exist %d", deploymentID)
	}
	deployment.Status = status
	deployment.EndTime = time.Now()
	return nil
}

// StoreEvent saves the events to the corresponding deploymentID
func (s *MemoryStorage) StoreEvent(event modelsdb.Event) error {
	deployment, exists := s.deployments.Get(event.ObjectID)
	if !exists {
		return fmt.Errorf("deployment doesn't exist %d", event.ObjectID)
	}
	deployment.Events = append(deployment.Events, &event)
	return nil
}

// GetEvents retreives all events related to the deploymentID
func (s *MemoryStorage) GetEvents(objectID uint64) ([]*modelsdb.Event, error) {
	deployment, exists := s.deployments.Get(objectID)
	if !exists {
		return nil, fmt.Errorf("deployment doesn't exist %d", objectID)
	}
	return deployment.Events, nil
}
