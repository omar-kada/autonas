package users

import (
	"testing"

	"omar-kada/autonas/internal/storage"
	"omar-kada/autonas/models"
	"omar-kada/autonas/testutil"

	"github.com/stretchr/testify/assert"
)

func TestLogin_Success(t *testing.T) {
	store, err := storage.NewUsersStorage(testutil.NewMemoryStorage())
	assert.NoError(t, err)
	service := NewService(store)

	credentials := models.Credentials{
		Username: "testuser",
		Password: "password",
	}

	hashedPassword, _ := hashPassword(credentials.Password)
	mockUser := models.User{
		Username:       credentials.Username,
		HashedPassword: hashedPassword,
	}

	store.UpsertUser(mockUser)

	token, err := service.Login(credentials)

	assert.NoError(t, err)
	assert.NotNil(t, token)
}

func TestLogin_UserNotFound(t *testing.T) {
	store, err := storage.NewUsersStorage(testutil.NewMemoryStorage())
	assert.NoError(t, err)
	service := NewService(store)

	credentials := models.Credentials{
		Username: "nonexistentuser",
		Password: "password",
	}

	token, err := service.Login(credentials)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrUserNotFound)
	assert.Zero(t, token)
}

func TestLogin_InvalidPassword(t *testing.T) {
	store, err := storage.NewUsersStorage(testutil.NewMemoryStorage())
	assert.NoError(t, err)
	service := NewService(store)

	credentials := models.Credentials{
		Username: "testuser",
		Password: "wrongpassword",
	}

	hashedPassword, _ := hashPassword("correctpassword")
	mockUser := models.User{
		Username:       credentials.Username,
		HashedPassword: hashedPassword,
	}

	store.UpsertUser(mockUser)

	token, err := service.Login(credentials)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidPassword)
	assert.Zero(t, token)
}

func TestIsRegistered_HasUsers(t *testing.T) {
	store, err := storage.NewUsersStorage(testutil.NewMemoryStorage())
	assert.NoError(t, err)
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
	store, err := storage.NewUsersStorage(testutil.NewMemoryStorage())
	assert.NoError(t, err)
	service := NewService(store)

	isRegistered, err := service.IsRegistered()

	assert.NoError(t, err)
	assert.False(t, isRegistered)
}

func TestRegister_Success(t *testing.T) {
	store, err := storage.NewUsersStorage(testutil.NewMemoryStorage())
	assert.NoError(t, err)
	service := NewService(store)

	credentials := models.Credentials{
		Username: "newuser",
		Password: "password",
	}

	token, err := service.Register(credentials)

	assert.NoError(t, err)
	assert.NotZero(t, token)
}

func TestRegister_AlreadyRegistered(t *testing.T) {
	store, err := storage.NewUsersStorage(testutil.NewMemoryStorage())
	assert.NoError(t, err)
	service := NewService(store)

	credentials := models.Credentials{
		Username: "existinguser",
		Password: "password",
	}

	mockUser := models.User{
		Username: credentials.Username,
	}

	store.UpsertUser(mockUser)

	token, err := service.Register(credentials)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrAlreadyRegistered)
	assert.Zero(t, token)
}

func TestLogout_Success(t *testing.T) {
	store, err := storage.NewUsersStorage(testutil.NewMemoryStorage())
	assert.NoError(t, err)
	service := NewService(store)

	pass, _ := hashPassword("password")
	mockUser := models.User{
		Username:       "testuser",
		HashedPassword: pass,
	}

	store.UpsertUser(mockUser)

	token, err := service.Login(models.Credentials{
		Username: "testuser",
		Password: "password",
	})
	assert.NoError(t, err)

	err = service.Logout(string(token.Value))
	assert.NoError(t, err)

	user, _ := service.GetUserByToken(string(token.Value))
	assert.Zero(t, user)
}

func TestLogout_UserNotFound(t *testing.T) {
	store, err := storage.NewUsersStorage(testutil.NewMemoryStorage())
	assert.NoError(t, err)
	service := NewService(store)

	err = service.Logout("invalidtoken")

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrUserNotFound)
}

func TestGetUserByToken_Success(t *testing.T) {
	store, err := storage.NewUsersStorage(testutil.NewMemoryStorage())
	assert.NoError(t, err)
	service := NewService(store)

	pass, _ := hashPassword("password")
	mockUser := models.User{
		Username:       "testuser",
		HashedPassword: pass,
	}

	store.UpsertUser(mockUser)
	token, err := service.Login(models.Credentials{
		Username: "testuser",
		Password: "password",
	})
	assert.NoError(t, err)
	user, err := service.GetUserByToken(string(token.Value))

	assert.NoError(t, err)
	assert.Equal(t, mockUser.Username, user.Username)
}

func TestGetUserByToken_UserNotFound(t *testing.T) {
	store, err := storage.NewUsersStorage(testutil.NewMemoryStorage())
	assert.NoError(t, err)
	service := NewService(store)

	user, err := service.GetUserByToken("invalidtoken")

	assert.ErrorIs(t, err, ErrUserNotFound)
	assert.Zero(t, user)
}

func TestGetUser_Success(t *testing.T) {
	store, err := storage.NewUsersStorage(testutil.NewMemoryStorage())
	assert.NoError(t, err)
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
	store, err := storage.NewUsersStorage(testutil.NewMemoryStorage())
	assert.NoError(t, err)
	service := NewService(store)

	username := "nonexistentuser"

	user, err := service.GetUser(username)

	assert.ErrorIs(t, err, ErrUserNotFound)
	assert.Empty(t, user.Username)
}

func TestDeleteUser_Success(t *testing.T) {
	store, err := storage.NewUsersStorage(testutil.NewMemoryStorage())
	assert.NoError(t, err)
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
	store, err := storage.NewUsersStorage(testutil.NewMemoryStorage())
	assert.NoError(t, err)
	service := NewService(store)

	username := "nonexistentuser"

	deleted, err := service.DeleteUser(username)

	assert.NoError(t, err)
	assert.False(t, deleted)
}

func TestChangePassword_Success(t *testing.T) {
	store, err := storage.NewUsersStorage(testutil.NewMemoryStorage())
	assert.NoError(t, err)
	service := NewService(store)

	username := "testuser"
	oldPassword := "oldpassword"
	newPassword := "newpassword"

	hashedPassword, _ := hashPassword(oldPassword)
	mockUser := models.User{
		Username:       username,
		HashedPassword: hashedPassword,
	}
	store.UpsertUser(mockUser)

	success, err := service.ChangePassword(username, oldPassword, newPassword)

	assert.NoError(t, err)
	assert.True(t, success)

	updatedUser, _ := store.UserByUsername(username)
	assert.True(t, checkPasswordHash(newPassword, updatedUser.HashedPassword))
}

func TestChangePassword_UserNotFound(t *testing.T) {
	store, err := storage.NewUsersStorage(testutil.NewMemoryStorage())
	assert.NoError(t, err)
	service := NewService(store)

	username := "nonexistentuser"
	oldPassword := "oldpassword"
	newPassword := "newpassword"

	success, err := service.ChangePassword(username, oldPassword, newPassword)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrUserNotFound)
	assert.False(t, success)
}

func TestChangePassword_InvalidOldPassword(t *testing.T) {
	store, err := storage.NewUsersStorage(testutil.NewMemoryStorage())
	assert.NoError(t, err)
	service := NewService(store)

	username := "testuser"
	oldPassword := "oldpassword"
	newPassword := "newpassword"
	wrongOldPassword := "wrongpassword"

	hashedPassword, _ := hashPassword(oldPassword)
	mockUser := models.User{
		Username:       username,
		HashedPassword: hashedPassword,
	}
	store.UpsertUser(mockUser)

	success, err := service.ChangePassword(username, wrongOldPassword, newPassword)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidPassword)
	assert.False(t, success)
}
