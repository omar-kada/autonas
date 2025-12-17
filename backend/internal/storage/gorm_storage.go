package storage

import (
	"errors"
	"time"

	"omar-kada/autonas/modelsdb"

	"gorm.io/gorm"
)

// gormStorage implements the Storage interface using GORM
type gormStorage struct {
	db *gorm.DB
}

// NewGormStorage creates a new instance of GormStorage and runs migrations
func NewGormStorage(db *gorm.DB) (Storage, error) {
	// Auto-migrate modelsdb types
	if err := db.AutoMigrate(&modelsdb.Deployment{}, &modelsdb.FileDiff{}, &modelsdb.Event{}); err != nil {
		return nil, err
	}
	return &gormStorage{db: db}, nil
}

// GetDeployments retrieves all deployments with their associated files and events
func (s *gormStorage) GetDeployments() ([]modelsdb.Deployment, error) {
	var deps []modelsdb.Deployment
	if err := s.db.Preload("Files").Preload("Events").Find(&deps).Error; err != nil {
		return nil, err
	}
	return deps, nil
}

// GetDeployment retrieves a specific deployment by ID with its associated files and events
func (s *gormStorage) GetDeployment(id uint64) (modelsdb.Deployment, error) {
	var dep modelsdb.Deployment
	if err := s.db.Preload("Files").Preload("Events").First(&dep, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return modelsdb.Deployment{}, err
		}
		return modelsdb.Deployment{}, err
	}
	return dep, nil
}

// InitDeployment creates a new deployment with the given parameters
func (s *gormStorage) InitDeployment(title string, author string, diff string, files []modelsdb.FileDiff) (modelsdb.Deployment, error) {
	dep := modelsdb.Deployment{
		Title:  title,
		Author: author,
		Diff:   diff,
		Status: modelsdb.DeploymentStatusRunning,
		Time:   time.Now(),
		Files:  files,
		Events: []modelsdb.Event{},
	}
	if err := s.db.Create(&dep).Error; err != nil {
		return modelsdb.Deployment{}, err
	}
	return dep, nil
}

// EndDeployment updates a deployment's status and sets its end time
func (s *gormStorage) EndDeployment(deploymentID uint64, status modelsdb.DeploymentStatus) error {
	var dep modelsdb.Deployment
	if err := s.db.First(&dep, deploymentID).Error; err != nil {
		return err
	}
	dep.Status = status
	dep.EndTime = time.Now()
	return s.db.Save(&dep).Error
}

// StoreEvent creates a new event and associates it with an existing deployment
func (s *gormStorage) StoreEvent(event modelsdb.Event) error {
	// verify deployment exists
	var dep modelsdb.Deployment
	if err := s.db.First(&dep, event.ObjectID).Error; err != nil {
		return err
	}
	if err := s.db.Create(&event).Error; err != nil {
		return err
	}
	return nil
}

// GetEvents retrieves all events associated with a specific object ID
func (s *gormStorage) GetEvents(objectID uint64) ([]modelsdb.Event, error) {
	var event []modelsdb.Event
	if err := s.db.Where("object_id = ?", objectID).Find(&event).Error; err != nil {
		return nil, err
	}
	return event, nil
}
