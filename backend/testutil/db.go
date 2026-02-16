package testutil

import (
	"log/slog"
	"testing"

	"omar-kada/autonas/internal/storage"

	"gorm.io/gorm"
)

// NewMemoryStorage instanciates a new memory storage
func NewMemoryStorage() *gorm.DB {
	db, err := storage.NewGormDb(":memory:", 0o000)
	if err != nil {
		slog.Error("couldn't init memory store")
	}
	return db
}

func NewDeploymentStorage(t *testing.T) storage.DeploymentStorage {
	depStore, err := storage.NewDeploymentStorage(NewMemoryStorage())
	if err != nil {
		t.Fatalf("error creating deployment storage : %v", err)
	}
	return depStore
}
