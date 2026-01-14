package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"omar-kada/autonas/models"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupStorage(t *testing.T) (Storage, *gorm.DB) {
	st, err := NewGormStorage(":memory:", 0o000)
	if err != nil {
		t.Fatalf("new storage: %v", err)
	}
	return st, st.(*gormStorage).db
}

func TestNewGormStorage_Migrates(t *testing.T) {
	_, db := setupStorage(t)
	// ensure migrations created the deployments table
	has := db.Migrator().HasTable(&models.Deployment{})
	assert.True(t, has)
}

func TestInitAndGetDeployment(t *testing.T) {
	s, _ := setupStorage(t)
	files := []models.FileDiff{{Diff: "d1", NewFile: "n1", OldFile: "o1"}}
	dep, err := s.InitDeployment("title1", "author1", "diff1", files)
	assert.NoError(t, err)
	assert.NotZero(t, dep.ID)
	assert.Equal(t, models.DeploymentStatusRunning, dep.Status)

	got, err := s.GetDeployment(dep.ID)
	assert.NoError(t, err)
	assert.Equal(t, dep.ID, got.ID)
	assert.Equal(t, files[0].Diff, got.Files[0].Diff)
	assert.Equal(t, files[0].NewFile, got.Files[0].NewFile)
	assert.Equal(t, files[0].OldFile, got.Files[0].OldFile)
	assert.Empty(t, got.Events)
}

func TestGetDeployment_NoNExisting(t *testing.T) {
	s, _ := setupStorage(t)

	dep, err := s.GetDeployment(999999)

	// Verify that no error is returned and the deployment is empty
	assert.NoError(t, err)
	assert.Equal(t, models.Deployment{}, dep)
}

func TestStoreEventAndGetEvents(t *testing.T) {
	s, _ := setupStorage(t)
	dep, err := s.InitDeployment("title1", "author1", "diff1", nil)
	assert.NoError(t, err)
	ev := models.Event{Level: 1, Msg: "ok", ObjectID: dep.ID}
	assert.NoError(t, s.StoreEvent(ev))
	events, err := s.GetEvents(dep.ID)
	assert.NoError(t, err)
	assert.Len(t, events, 1)
	assert.Equal(t, "ok", events[0].Msg)
}

func TestGetLastAndGetDeploymentsOrdering(t *testing.T) {
	s, _ := setupStorage(t)
	_, err := s.InitDeployment("title1", "author1", "diff1", nil)
	assert.NoError(t, err)

	// small sleep to ensure time difference
	time.Sleep(2 * time.Millisecond)
	dep2, _ := s.InitDeployment("title2", "author2", "diff2", nil)
	time.Sleep(2 * time.Millisecond)
	dep3, _ := s.InitDeployment("title3", "author3", "diff3", nil)

	last, err := s.GetLastDeployment()
	assert.NoError(t, err)
	assert.Equal(t, dep3.ID, last.ID)

	deps, err := s.GetDeployments(Cursor[uint64]{Limit: 10, Offset: 999999})
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(deps), 3)
	assert.Equal(t, dep3.ID, deps[0].ID)
	assert.Equal(t, dep2.ID, deps[1].ID)
}

func TestEndDeploymentAndErrorCases(t *testing.T) {
	s, _ := setupStorage(t)
	dep, err := s.InitDeployment("title1", "author1", "diff1", nil)
	assert.NoError(t, err)

	assert.NoError(t, s.EndDeployment(dep.ID, models.DeploymentStatusSuccess))
	d, err := s.GetDeployment(dep.ID)
	assert.NoError(t, err)
	assert.Equal(t, models.DeploymentStatusSuccess, d.Status)
	assert.False(t, d.EndTime.IsZero())

	// StoreEvent should fail for a non-existing deployment
	err = s.StoreEvent(models.Event{Level: 1, Msg: "bad", ObjectID: 999999})
	assert.Error(t, err)
}

func TestGetDeployments_Pagination(t *testing.T) {
	cases := []struct {
		name     string
		seed     int
		cursor   Cursor[uint64]
		expected int
	}{
		{"DefaultOffset", 15, NewIDCursor(5, 0), 5},
		{"LimitCapped", 50, NewIDCursor(200, 0), 50},
		{"OffsetWithinRange", 15, NewIDCursor(10, 10), 9},

		{"ExactLimit", 10, NewIDCursor(3, 0), 3},
		{"NoResults", 10, Cursor[uint64]{Limit: 10, Offset: 1}, 0},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s, _ := setupStorage(t)
			for i := 1; i <= c.seed; i++ {
				_, err := s.InitDeployment(fmt.Sprintf("t%d", i), "author", "diff", nil)
				assert.NoError(t, err)
			}

			deps, err := s.GetDeployments(c.cursor)
			assert.NoError(t, err)
			assert.Equal(t, c.expected, len(deps))

			for _, d := range deps {
				assert.Less(t, d.ID, c.cursor.Offset)
			}
		})
	}
}

func TestNewGormStorage_FileCreation(t *testing.T) {
	// Create a temporary directory for the test
	tempDir := t.TempDir()
	dbFile := filepath.Join(tempDir, "test.db")

	// Create a new GORM storage with the temporary file
	st, err := NewGormStorage(dbFile, 0o000)
	assert.NoError(t, err)
	assert.NotNil(t, st)

	// Verify that the database file was created
	_, err = os.Stat(dbFile)
	assert.NoError(t, err, "Database file should be created")

	// Clean up: close the database connection
	db, ok := st.(*gormStorage)
	assert.True(t, ok, "Expected gormStorage type")
	sqlDB, err := db.db.DB()
	assert.NoError(t, err)
	assert.NoError(t, sqlDB.Close())
}
