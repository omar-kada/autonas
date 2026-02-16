// Package users provides user management and authentication services.
package users

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"omar-kada/autonas/internal/storage"
	"omar-kada/autonas/models"

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

// use in memory storage for token , and user a timer to remove it from the map once expired,
//

var usersToken = make(map[TokenValue]string)

// Service abstracts authorization operations
type Service interface {
	AuthService
	AccountService
}

// AuthService abstracts authentication operations
type AuthService interface {
	Login(credentials models.Credentials) (Token, error)
	Register(credentials models.Credentials) (Token, error)
	IsRegistered() (bool, error)
	Logout(token string) error
	GetUserByToken(token string) (models.User, error)
}

// AccountService abstracts account management operations
type AccountService interface {
	GetUser(username string) (models.User, error)
	DeleteUser(username string) (bool, error)
	ChangePassword(username string, oldPass string, newPass string) (bool, error)
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
func (a *service) Login(credentials models.Credentials) (Token, error) {
	user, err := a.userStore.UserByUsername(credentials.Username)
	if err != nil {
		return Token{}, fmt.Errorf("error finding user: %w", err)
	}

	if user.Username == "" {
		return Token{}, ErrUserNotFound
	}

	if !checkPasswordHash(credentials.Password, user.HashedPassword) {
		return Token{}, ErrInvalidPassword
	}

	token, err := generateToken()
	if err != nil {
		return Token{}, fmt.Errorf("error generating token: %w", err)
	}

	a.insertToken(token, user)

	return token, nil
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
func (a *service) Register(credentials models.Credentials) (Token, error) {
	hasUsers, err := a.userStore.HasUsers()
	if err != nil {
		return Token{}, err
	}
	if hasUsers {
		return Token{}, ErrAlreadyRegistered
	}
	hashedPassword, err := hashPassword(credentials.Password)
	if err != nil {
		return Token{}, fmt.Errorf("error hashing password: %w", err)
	}

	token, err := generateToken()
	if err != nil {
		return Token{}, fmt.Errorf("error generating token: %w", err)
	}

	user := models.User{
		Username:       credentials.Username,
		HashedPassword: hashedPassword,
	}
	_, err = a.userStore.UpsertUser(user)
	if err != nil {
		return Token{}, fmt.Errorf("error creating user: %w", err)
	}

	a.insertToken(token, user)

	return token, nil
}

func (*service) insertToken(token Token, user models.User) {
	usersToken[token.Value] = user.Username
	go func() {
		time.AfterFunc(time.Until(token.Expires), func() {
			delete(usersToken, token.Value)
		})
	}()
}

// Logout invalidates the user's authentication token.
func (a *service) Logout(tokenValue string) error {
	if _, exists := usersToken[TokenValue(tokenValue)]; !exists {
		return ErrUserNotFound
	}

	delete(usersToken, TokenValue(tokenValue))
	return nil
}

// GetUserByToken retrieves a user by their authentication token.
func (a *service) GetUserByToken(tokenValue string) (models.User, error) {
	user, err := a.userStore.UserByUsername(usersToken[TokenValue(tokenValue)])
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

// ChangePassword updates a user's password after verifying the old password.
func (a *service) ChangePassword(username string, oldPass string, newPass string) (bool, error) {
	user, err := a.userStore.UserByUsername(username)
	if err != nil {
		return false, fmt.Errorf("error finding user: %w", err)
	}

	if user.Username == "" {
		return false, ErrUserNotFound
	}

	if !checkPasswordHash(oldPass, user.HashedPassword) {
		return false, ErrInvalidPassword
	}

	hashedPassword, err := hashPassword(newPass)
	if err != nil {
		return false, fmt.Errorf("error hashing new password: %w", err)
	}

	user.HashedPassword = hashedPassword
	if _, err := a.userStore.UpsertUser(user); err != nil {
		return false, fmt.Errorf("error updating user: %w", err)
	}

	return true, nil
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

type (
	TokenValue string
	Token      struct {
		Value   TokenValue
		Expires time.Time
	}
)

func generateToken() (Token, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return Token{}, err
	}
	return Token{
		Value:   TokenValue(hex.EncodeToString(token)),
		Expires: time.Now().Add(30 * time.Minute),
	}, nil
}
