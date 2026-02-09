package storage

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"omar-kada/autonas/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

var ErrEmptyUsername = errors.New("empty username")

// gormStorage implements the Storage interface using GORM
type gormStorage struct {
	db *gorm.DB
}

// NewGormStorage creates a new instance of GormStorage and runs migrations
func NewGormStorage(dbFile string, addPerm os.FileMode) (Storage, error) {
	if dbFile != ":memory:" {
		if _, err := os.Stat(dbFile); os.IsNotExist(err) {
			if err := os.MkdirAll(filepath.Dir(dbFile), 0o700|addPerm); err != nil {
				return nil, err
			}
		} else if err != nil {
			return nil, err
		}
	}
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Warn, // Log level
			Colorful:      true,        // Enable color
		},
	)

	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}
	// pragmas and pooling
	db.Exec("PRAGMA journal_mode=WAL;")
	db.Exec("PRAGMA foreign_keys = ON;")
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("couldn't init sqlite db %w", err)
	}
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)

	// Auto-migrate models types
	if err := db.AutoMigrate(&models.Deployment{}, &models.FileDiff{}, &models.Event{}, &models.User{}); err != nil {
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
			return models.Deployment{}, nil
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

// HasUsers checks if there are any users in the database
func (s *gormStorage) HasUsers() (bool, error) {
	var exists bool
	if err := s.db.Model(&models.User{}).Select("EXISTS (SELECT 1 FROM users)").Find(&exists).Error; err != nil {
		return false, err
	}
	return exists, nil
}

// UserByToken retrieves a user by their token
func (s *gormStorage) UserByToken(token string) (models.User, error) {
	var user models.User
	if err := s.db.Where("token = ?", token).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, nil
		}
		return models.User{}, err
	}
	return user, nil
}

// UserByUsername retrieves a user by their username
func (s *gormStorage) UserByUsername(username string) (models.User, error) {
	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, nil
		}
		return models.User{}, err
	}
	return user, nil
}

// UpsertUser updates an existing user in the database
func (s *gormStorage) UpsertUser(user models.User) (models.User, error) {
	if user.Username == "" {
		return models.User{}, ErrEmptyUsername
	}
	if err := s.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (s *gormStorage) DeleteUserByUserName(username string) (bool, error) {
	tx := s.db.Where("username = ?", username).Delete(&models.User{})
	if err := tx.Error; err != nil {
		return false, err
	}
	return tx.RowsAffected > 0, nil
}
