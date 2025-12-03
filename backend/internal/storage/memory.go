package storage

import (
	"fmt"
	"log/slog"
	"omar-kada/autonas/api"
)

// MemoryStorage uses memory to store data (to be used mainly for testing)
type MemoryStorage struct {
	deployments map[string]api.Deployment
}

// NewMemoryStorage instanciates a new memory storage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		deployments: make(map[string]api.Deployment),
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
		deployments = append(deployments, deployment)
	}

	return deployments, nil
}

// SaveDeployment saves a deployment
func (s *MemoryStorage) SaveDeployment(deployment api.Deployment) (api.Deployment, error) {
	logs := deployment.Logs
	if existing, exists := s.deployments[deployment.Id]; exists {
		logs = append(existing.Logs, logs...)
	}
	deployment.Logs = logs
	s.deployments[deployment.Id] = deployment
	return deployment, nil
}

// UpdateStatus updates only the status of the deployment
func (s *MemoryStorage) UpdateStatus(deploymentID string, status api.DeploymentStatus) (api.Deployment, error) {
	deployment, exists := s.deployments[deploymentID]
	if !exists {
		return api.Deployment{}, nil
	}
	deployment.Status = status
	return deployment, nil
}

// AddLogRecord adds a log record to a deployment
func (s *MemoryStorage) AddLogRecord(deploymentID string, record slog.Record) error {
	deployment, exists := s.deployments[deploymentID]
	if !exists {
		fmt.Println("deployment doesn't exist", deploymentID, record)
		return nil
	}
	fmt.Println("adding record to deployment", deploymentID, record)

	deployment.Logs = append(deployment.Logs, fmt.Sprintf("%v: [%s] %s", record.Time.Format("15:04:05"), record.Level.String(), record.Message))
	return nil
}
