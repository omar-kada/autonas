// Package user provides user management and authentication services.
package user

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"omar-kada/autonas/internal/storage"
	"omar-kada/autonas/models"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrAlreadyRegistered indicates that a user is already registered
	ErrAlreadyRegistered = errors.New("already registered")
	// ErrUserNotFound indicates that a user was not found
	ErrUserNotFound = errors.New("user not found")
	// ErrInvalidPassword indicates that the provided password is incorrect
	ErrInvalidPassword = errors.New("invalid password")
)

// Service abstracts authorization operations
type Service interface {
	AuthService
	AccountService
}

// AuthService abstracts authentication operations
type AuthService interface {
	Login(credentials models.Credentials) (models.Auth, error)
	Register(credentials models.Credentials) (models.Auth, error)
	IsRegistered() (bool, error)
	Logout(token string) error
	GetUserByToken(token string) (models.User, error)
}

// AccountService abstracts account management operations
type AccountService interface {
	GetUser(username string) (models.User, error)
	DeleteUser(username string) (bool, error)
}

type service struct {
	userStore storage.UserStorage
}

// NewService creates a new userService
func NewService(userStore storage.UserStorage) Service {
	return &service{
		userStore: userStore,
	}
}

// Login authenticates a user and returns their authentication token.
func (a *service) Login(credentials models.Credentials) (models.Auth, error) {
	user, err := a.userStore.UserByUsername(credentials.Username)
	if err != nil {
		return models.Auth{}, fmt.Errorf("error finding user: %w", err)
	}

	if user.Username == "" {
		return models.Auth{}, ErrUserNotFound
	}

	if !checkPasswordHash(credentials.Password, user.HasedPassword) {
		return models.Auth{}, ErrInvalidPassword
	}

	token, err := generateToken()
	if err != nil {
		return models.Auth{}, fmt.Errorf("error generating token: %w", err)
	}

	user.Token = token
	user.ExpiresIn = time.Now().Add(24 * time.Hour)
	savedUsed, err := a.userStore.UpsertUser(user)
	if err != nil {
		return models.Auth{}, fmt.Errorf("error updating user: %w", err)
	}

	return savedUsed.Auth, nil
}

// IsRegistered checks if any users are registered in the system.
func (a *service) IsRegistered() (bool, error) {
	hasUsers, err := a.userStore.HasUsers()
	if err != nil {
		return false, fmt.Errorf("error checking if users exist: %w", err)
	}
	return hasUsers, nil
}

// Register creates a new user account with the provided credentials
func (a *service) Register(credentials models.Credentials) (models.Auth, error) {
	hasUsers, err := a.userStore.HasUsers()
	if err != nil {
		return models.Auth{}, err
	}
	if hasUsers {
		return models.Auth{}, ErrAlreadyRegistered
	}
	hashedPassword, err := hashPassword(credentials.Password)
	if err != nil {
		return models.Auth{}, fmt.Errorf("error hashing password: %w", err)
	}

	token, err := generateToken()
	if err != nil {
		return models.Auth{}, fmt.Errorf("error generating token: %w", err)
	}

	user := models.User{
		Username:      credentials.Username,
		HasedPassword: hashedPassword,
		Auth: models.Auth{
			Token:     token,
			ExpiresIn: time.Now().Add(24 * time.Hour),
		},
	}
	savedUser, err := a.userStore.UpsertUser(user)
	if err != nil {
		return models.Auth{}, fmt.Errorf("error creating user: %w", err)
	}

	return savedUser.Auth, nil
}

// Logout invalidates the user's authentication token.
func (a *service) Logout(token string) error {
	user, err := a.userStore.UserByToken(token)
	if err != nil {
		return fmt.Errorf("error finding user: %w", err)
	}

	if user.Username == "" {
		return ErrUserNotFound
	}

	user.Token = ""
	user.ExpiresIn = time.Time{}

	if _, err := a.userStore.UpsertUser(user); err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	return nil
}

// GetUserByToken retrieves a user by their authentication token.
func (a *service) GetUserByToken(token string) (models.User, error) {
	user, err := a.userStore.UserByToken(token)
	if err == nil && user.Username == "" {
		return user, ErrUserNotFound
	}
	return user, err
}

// GetUser retrieves a user by their username.
func (a *service) GetUser(username string) (models.User, error) {
	user, err := a.userStore.UserByUsername(username)
	if err == nil && user.Username == "" {
		return user, ErrUserNotFound
	}
	return user, err
}

// DeleteUser removes a user from the system by their username.
func (a *service) DeleteUser(username string) (bool, error) {
	return a.userStore.DeleteUserByUserName(username)
}

// Helper functions
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}
