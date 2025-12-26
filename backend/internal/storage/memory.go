package storage

import (
	"fmt"
	"sync"
	"time"

	"omar-kada/autonas/models"

	"github.com/elliotchance/orderedmap/v3"
)

// MemoryStorage uses memory to store data (to be used mainly for testing)
type MemoryStorage struct {
	deployments *orderedmap.OrderedMap[uint64, models.Deployment]
	mu          sync.Mutex
}

// NewMemoryStorage instanciates a new memory storage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		deployments: orderedmap.NewOrderedMap[uint64, models.Deployment](),
	}
}

// GetDeployments returns stored deployments
func (s *MemoryStorage) GetDeployments(c Cursor[uint64]) ([]models.Deployment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// transform map to slice
	var deployments []models.Deployment
	startIndex := 0
	for i, deployment := range s.deployments.AllFromBack() {
		deployments = append(deployments, deployment)
		if deployment.ID == c.Offset {
			startIndex = int(i)
		}
	}
	limit := min(startIndex+c.Limit, len(deployments))

	return deployments[startIndex:limit], nil
}

// GetDeployment returns deployment by id
func (s *MemoryStorage) GetDeployment(id uint64) (models.Deployment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	deployment, exists := s.deployments.Get(id)
	if !exists {
		return models.Deployment{}, fmt.Errorf("deployment doesn't exist %d", id)
	}
	return deployment, nil
}

// GetLastDeployment returns the most recent deployment or an error if none exist
func (s *MemoryStorage) GetLastDeployment() (models.Deployment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, deployment := range s.deployments.AllFromBack() {
		return deployment, nil
	}
	return models.Deployment{}, fmt.Errorf("no deployments available")
}

func newID() uint64 {
	id := uint64(time.Now().UnixNano())
	return id
}

// InitDeployment creates a new deployment and returns it
func (s *MemoryStorage) InitDeployment(title string, author string, diff string, files []models.FileDiff) (models.Deployment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	deployment := models.Deployment{
		ID:     newID(),
		Title:  title,
		Author: author,
		Time:   time.Now(),
		Status: "running",
		Diff:   diff,
		Files:  files,
		Events: []models.Event{},
	}
	s.deployments.Set(deployment.ID, deployment)
	return deployment, nil
}

// EndDeployment updates only the status of the deployment
func (s *MemoryStorage) EndDeployment(deploymentID uint64, status models.DeploymentStatus) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	deployment, exists := s.deployments.Get(deploymentID)
	if !exists {
		return fmt.Errorf("deployment doesn't exist %d", deploymentID)
	}
	deployment.Status = status
	deployment.EndTime = time.Now()
	s.deployments.Set(deploymentID, deployment)
	return nil
}

// StoreEvent saves the events to the corresponding deploymentID
func (s *MemoryStorage) StoreEvent(event models.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	deployment, exists := s.deployments.Get(event.ObjectID)
	if !exists {
		return fmt.Errorf("deployment doesn't exist %d", event.ObjectID)
	}
	deployment.Events = append(deployment.Events, event)
	s.deployments.Set(event.ObjectID, deployment)

	return nil
}

// GetEvents retreives all events related to the deploymentID
func (s *MemoryStorage) GetEvents(objectID uint64) ([]models.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	deployment, exists := s.deployments.Get(objectID)
	if !exists {
		return nil, fmt.Errorf("deployment doesn't exist %d", objectID)
	}
	return deployment.Events, nil
}
