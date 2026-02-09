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

func TestHasUsers(t *testing.T) {
	s, _ := setupStorage(t)

	// Test when there are no users
	exists, err := s.HasUsers()
	assert.NoError(t, err)
	assert.False(t, exists)

	// Add a user to the database
	user := models.User{
		Username:       "testuser",
		HashedPassword: "hashedpassword",
		Auth: models.Auth{
			Token:     "testtoken",
			ExpiresIn: time.Now().Add(24 * time.Hour),
		},
	}
	_, err = s.UpsertUser(user)
	assert.NoError(t, err)

	// Test when there are users
	exists, err = s.HasUsers()
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestUserByToken(t *testing.T) {
	s, _ := setupStorage(t)

	// Test when the token does not exist
	user, err := s.UserByToken("nonexistenttoken")
	assert.NoError(t, err)
	assert.Equal(t, models.User{}, user)

	// Add a user to the database
	userToAdd := models.User{
		Username:       "testuser",
		HashedPassword: "hashedpassword",
		Auth: models.Auth{
			Token:     "testtoken",
			ExpiresIn: time.Now().Add(24 * time.Hour),
		},
	}
	_, err = s.UpsertUser(userToAdd)
	assert.NoError(t, err)

	// Test when the token exists
	user, err = s.UserByToken("testtoken")
	assert.NoError(t, err)
	assert.Equal(t, userToAdd.Username, user.Username)
	assert.Equal(t, userToAdd.HashedPassword, user.HashedPassword)
	assert.Equal(t, userToAdd.Token, user.Token)
	assert.Equal(t, userToAdd.ExpiresIn.Unix(), user.ExpiresIn.Unix())
}
func TestUserByUsername(t *testing.T) {
	s, _ := setupStorage(t)

	// Test when the username does not exist
	user, err := s.UserByUsername("nonexistentuser")
	assert.NoError(t, err)
	assert.Equal(t, models.User{}, user)

	// Add a user to the database
	userToAdd := models.User{
		Username:       "testuser",
		HashedPassword: "hashedpassword",
		Auth: models.Auth{
			Token:     "testtoken",
			ExpiresIn: time.Now().Add(24 * time.Hour),
		},
	}
	_, err = s.UpsertUser(userToAdd)
	assert.NoError(t, err)

	// Test when the username exists
	user, err = s.UserByUsername("testuser")
	assert.NoError(t, err)
	assert.Equal(t, userToAdd.Username, user.Username)
	assert.Equal(t, userToAdd.HashedPassword, user.HashedPassword)
	assert.Equal(t, userToAdd.Token, user.Token)
	assert.Equal(t, userToAdd.ExpiresIn.Unix(), user.ExpiresIn.Unix())
}
func TestUpsertUser_InsertNewUser(t *testing.T) {
	s, _ := setupStorage(t)

	// Create a new user
	user := models.User{
		Username:       "newuser",
		HashedPassword: "newpassword",
		Auth: models.Auth{
			Token:     "newtoken",
			ExpiresIn: time.Now().Add(24 * time.Hour),
		},
	}

	// Insert the new user
	insertedUser, err := s.UpsertUser(user)
	assert.NoError(t, err)
	assert.Equal(t, user.Username, insertedUser.Username)
	assert.Equal(t, user.HashedPassword, insertedUser.HashedPassword)
	assert.Equal(t, user.Token, insertedUser.Token)
	assert.Equal(t, user.ExpiresIn.Unix(), insertedUser.ExpiresIn.Unix())

	// Verify the user was inserted
	fetchedUser, err := s.UserByUsername("newuser")
	assert.NoError(t, err)
	assert.Equal(t, user.Username, fetchedUser.Username)
	assert.Equal(t, user.HashedPassword, fetchedUser.HashedPassword)
	assert.Equal(t, user.Token, fetchedUser.Token)
	assert.Equal(t, user.ExpiresIn.Unix(), fetchedUser.ExpiresIn.Unix())
}

func TestUpsertUser_UpdateExistingUser(t *testing.T) {
	s, _ := setupStorage(t)

	// Insert an initial user
	initialUser := models.User{
		Username:       "existinguser",
		HashedPassword: "initialpassword",
		Auth: models.Auth{
			Token:     "initialtoken",
			ExpiresIn: time.Now().Add(24 * time.Hour),
		},
	}
	_, err := s.UpsertUser(initialUser)
	assert.NoError(t, err)

	// Update the user
	updatedUser := models.User{
		Username:       "existinguser",
		HashedPassword: "updatedpassword",
		Auth: models.Auth{
			Token:     "updatedtoken",
			ExpiresIn: time.Now().Add(48 * time.Hour),
		},
	}
	upsertedUser, err := s.UpsertUser(updatedUser)
	assert.NoError(t, err)
	assert.Equal(t, updatedUser.Username, upsertedUser.Username)
	assert.Equal(t, updatedUser.HashedPassword, upsertedUser.HashedPassword)
	assert.Equal(t, updatedUser.Token, upsertedUser.Token)
	assert.Equal(t, updatedUser.ExpiresIn.Unix(), upsertedUser.ExpiresIn.Unix())

	// Verify the user was updated
	fetchedUser, err := s.UserByUsername("existinguser")
	assert.NoError(t, err)
	assert.Equal(t, updatedUser.Username, fetchedUser.Username)
	assert.Equal(t, updatedUser.HashedPassword, fetchedUser.HashedPassword)
	assert.Equal(t, updatedUser.Token, fetchedUser.Token)
	assert.Equal(t, updatedUser.ExpiresIn.Unix(), fetchedUser.ExpiresIn.Unix())
}

func TestUpsertUser_InvalidUser(t *testing.T) {
	s, _ := setupStorage(t)

	// Attempt to insert a user with invalid data
	invalidUser := models.User{
		Username:       "",
		HashedPassword: "",
		Auth: models.Auth{
			Token:     "",
			ExpiresIn: time.Time{},
		},
	}
	_, err := s.UpsertUser(invalidUser)
	assert.ErrorIs(t, err, ErrEmptyUsername)
}
func TestDeleteUserByUserName_Success(t *testing.T) {
	s, _ := setupStorage(t)

	// Add a user to the database
	userToAdd := models.User{
		Username:       "testuser",
		HashedPassword: "hashedpassword",
		Auth: models.Auth{
			Token:     "testtoken",
			ExpiresIn: time.Now().Add(24 * time.Hour),
		},
	}
	_, err := s.UpsertUser(userToAdd)
	assert.NoError(t, err)

	// Delete the user by username
	deleted, err := s.DeleteUserByUserName("testuser")
	assert.NoError(t, err)
	assert.True(t, deleted)

	// Verify the user was deleted
	user, err := s.UserByUsername("testuser")
	assert.NoError(t, err)
	assert.Equal(t, models.User{}, user)
}

func TestDeleteUserByUserName_NotFound(t *testing.T) {
	s, _ := setupStorage(t)

	// Attempt to delete a non-existent user
	deleted, err := s.DeleteUserByUserName("nonexistentuser")
	assert.NoError(t, err)
	assert.False(t, deleted)
}
