package storage

import (
	"errors"
	"time"

	"omar-kada/autonas/models"

	"gorm.io/gorm"
)

// gormStorage implements the Storage interface using GORM
type gormStorage struct {
	db *gorm.DB
}

// NewGormStorage creates a new instance of GormStorage and runs migrations
func NewGormStorage(db *gorm.DB) (Storage, error) {
	// Auto-migrate models types
	if err := db.AutoMigrate(&models.Deployment{}, &models.FileDiff{}, &models.Event{}); err != nil {
		return nil, err
	}
	return &gormStorage{db: db}, nil
}

// GetDeployments retrieves all deployments with their associated files and events
func (s *gormStorage) GetDeployments(c Cursor[uint64]) ([]models.Deployment, error) {
	var deps []models.Deployment
	if err := s.db.
		Scopes(Paginate(c)).Order("Time desc").Find(&deps).Error; err != nil {
		return nil, err
	}
	return deps, nil
}

// GetDeployment retrieves a specific deployment by ID with its associated files and events
func (s *gormStorage) GetDeployment(id uint64) (models.Deployment, error) {
	var dep models.Deployment
	if err := s.db.Preload("Files").Preload("Events").First(&dep, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Deployment{}, err
		}
		return models.Deployment{}, err
	}
	return dep, nil
}

// InitDeployment creates a new deployment with the given parameters
func (s *gormStorage) InitDeployment(title string, author string, diff string, files []models.FileDiff) (models.Deployment, error) {
	dep := models.Deployment{
		Title:  title,
		Author: author,
		Diff:   diff,
		Status: models.DeploymentStatusRunning,
		Time:   time.Now(),
		Files:  files,
		Events: []models.Event{},
	}
	if err := s.db.Create(&dep).Error; err != nil {
		return models.Deployment{}, err
	}
	return dep, nil
}

// EndDeployment updates a deployment's status and sets its end time
func (s *gormStorage) EndDeployment(deploymentID uint64, status models.DeploymentStatus) error {
	var dep models.Deployment
	if err := s.db.First(&dep, deploymentID).Error; err != nil {
		return err
	}
	dep.Status = status
	dep.EndTime = time.Now()
	return s.db.Save(&dep).Error
}

// GetLastDeployment returns the most recent deployment based on Time (or ID) descending
func (s *gormStorage) GetLastDeployment() (models.Deployment, error) {
	var dep models.Deployment
	req := s.db.Preload("Files").Preload("Events").Order("time DESC")
	if err := req.First(&dep).Error; err != nil {
		return models.Deployment{}, err
	}
	return dep, nil
}

// StoreEvent creates a new event and associates it with an existing deployment
func (s *gormStorage) StoreEvent(event models.Event) error {
	// verify deployment exists
	var dep models.Deployment
	if err := s.db.First(&dep, event.ObjectID).Error; err != nil {
		return err
	}
	if err := s.db.Create(&event).Error; err != nil {
		return err
	}
	return nil
}

// GetEvents retrieves all events associated with a specific object ID
func (s *gormStorage) GetEvents(objectID uint64) ([]models.Event, error) {
	var event []models.Event
	if err := s.db.Where("object_id = ?", objectID).Find(&event).Error; err != nil {
		return nil, err
	}
	return event, nil
}
