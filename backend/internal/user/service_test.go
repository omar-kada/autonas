package user

import (
	"testing"
	"time"

	"omar-kada/autonas/models"
	"omar-kada/autonas/testutil"

	"github.com/stretchr/testify/assert"
)

func TestLogin_Success(t *testing.T) {
	store := testutil.NewMemoryStorage()
	service := NewService(store)

	credentials := models.Credentials{
		Username: "testuser",
		Password: "password",
	}

	hashedPassword, _ := hashPassword(credentials.Password)
	mockUser := models.User{
		Username:      credentials.Username,
		HashedPassword: hashedPassword,
	}

	store.UpsertUser(mockUser)

	auth, err := service.Login(credentials)

	assert.NoError(t, err)
	assert.NotEmpty(t, auth.Token)
	assert.NotZero(t, auth.ExpiresIn)
}

func TestLogin_UserNotFound(t *testing.T) {
	store := testutil.NewMemoryStorage()
	service := NewService(store)

	credentials := models.Credentials{
		Username: "nonexistentuser",
		Password: "password",
	}

	auth, err := service.Login(credentials)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrUserNotFound)
	assert.Empty(t, auth.Token)
}

func TestLogin_InvalidPassword(t *testing.T) {
	store := testutil.NewMemoryStorage()
	service := NewService(store)

	credentials := models.Credentials{
		Username: "testuser",
		Password: "wrongpassword",
	}

	hashedPassword, _ := hashPassword("correctpassword")
	mockUser := models.User{
		Username:      credentials.Username,
		HashedPassword: hashedPassword,
	}

	store.UpsertUser(mockUser)

	auth, err := service.Login(credentials)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidPassword)
	assert.Empty(t, auth.Token)
}

func TestIsRegistered_HasUsers(t *testing.T) {
	store := testutil.NewMemoryStorage()
	service := NewService(store)

	mockUser := models.User{
		Username: "testuser",
	}

	store.UpsertUser(mockUser)

	isRegistered, err := service.IsRegistered()

	assert.NoError(t, err)
	assert.True(t, isRegistered)
}

func TestIsRegistered_NoUsers(t *testing.T) {
	store := testutil.NewMemoryStorage()
	service := NewService(store)

	isRegistered, err := service.IsRegistered()

	assert.NoError(t, err)
	assert.False(t, isRegistered)
}

func TestRegister_Success(t *testing.T) {
	store := testutil.NewMemoryStorage()
	service := NewService(store)

	credentials := models.Credentials{
		Username: "newuser",
		Password: "password",
	}

	auth, err := service.Register(credentials)

	assert.NoError(t, err)
	assert.NotEmpty(t, auth.Token)
	assert.NotZero(t, auth.ExpiresIn)
}

func TestRegister_AlreadyRegistered(t *testing.T) {
	store := testutil.NewMemoryStorage()
	service := NewService(store)

	credentials := models.Credentials{
		Username: "existinguser",
		Password: "password",
	}

	mockUser := models.User{
		Username: credentials.Username,
	}

	store.UpsertUser(mockUser)

	auth, err := service.Register(credentials)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrAlreadyRegistered)
	assert.Empty(t, auth.Token)
}

func TestLogout_Success(t *testing.T) {
	store := testutil.NewMemoryStorage()
	service := NewService(store)

	token := "validtoken"
	mockUser := models.User{
		Username: "testuser",
		Auth:     models.Auth{Token: token, ExpiresIn: time.Now().Add(24 * time.Hour)},
	}

	store.UpsertUser(mockUser)

	err := service.Logout(token)

	assert.NoError(t, err)

	user, _ := store.UserByToken(token)
	assert.Empty(t, user.Token)
}

func TestLogout_UserNotFound(t *testing.T) {
	store := testutil.NewMemoryStorage()
	service := NewService(store)

	token := "invalidtoken"

	err := service.Logout(token)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrUserNotFound)
}

func TestGetUserByToken_Success(t *testing.T) {
	store := testutil.NewMemoryStorage()
	service := NewService(store)

	token := "validtoken"
	mockUser := models.User{
		Username: "testuser",
		Auth:     models.Auth{Token: token, ExpiresIn: time.Now().Add(24 * time.Hour)},
	}

	store.UpsertUser(mockUser)

	user, err := service.GetUserByToken(token)

	assert.NoError(t, err)
	assert.Equal(t, mockUser.Username, user.Username)
}

func TestGetUserByToken_UserNotFound(t *testing.T) {
	store := testutil.NewMemoryStorage()
	service := NewService(store)

	token := "invalidtoken"

	user, err := service.GetUserByToken(token)

	assert.ErrorIs(t, err, ErrUserNotFound)
	assert.Empty(t, user.Username)
}

func TestGetUser_Success(t *testing.T) {
	store := testutil.NewMemoryStorage()
	service := NewService(store)

	username := "testuser"
	mockUser := models.User{
		Username: username,
	}

	store.UpsertUser(mockUser)

	user, err := service.GetUser(username)

	assert.NoError(t, err)
	assert.Equal(t, mockUser, user)
}

func TestGetUser_UserNotFound(t *testing.T) {
	store := testutil.NewMemoryStorage()
	service := NewService(store)

	username := "nonexistentuser"

	user, err := service.GetUser(username)

	assert.ErrorIs(t, err, ErrUserNotFound)
	assert.Empty(t, user.Username)
}

func TestDeleteUser_Success(t *testing.T) {
	store := testutil.NewMemoryStorage()
	service := NewService(store)

	username := "testuser"
	mockUser := models.User{
		Username: username,
	}

	store.UpsertUser(mockUser)

	deleted, err := service.DeleteUser(username)

	assert.NoError(t, err)
	assert.True(t, deleted)
}

func TestDeleteUser_UserNotFound(t *testing.T) {
	store := testutil.NewMemoryStorage()
	service := NewService(store)

	username := "nonexistentuser"

	deleted, err := service.DeleteUser(username)

	assert.NoError(t, err)
	assert.False(t, deleted)
}
