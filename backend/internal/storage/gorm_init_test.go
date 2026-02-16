package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupStorage(t *testing.T) *gorm.DB {
	db, err := NewGormDb(":memory:", 0o000)
	if err != nil {
		t.Fatalf("new storage: %v", err)
	}
	return db
}

func TestNewGormStorage_FileCreation(t *testing.T) {
	// Create a temporary directory for the test
	tempDir := t.TempDir()
	dbFile := filepath.Join(tempDir, "test.db")

	// Create a new GORM storage with the temporary file
	db, err := NewGormDb(dbFile, 0o000)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	// Verify that the database file was created
	_, err = os.Stat(dbFile)
	assert.NoError(t, err, "Database file should be created")

	// Clean up: close the database connection
	sqlDB, err := db.DB()
	assert.NoError(t, err)
	assert.NoError(t, sqlDB.Close())
}
