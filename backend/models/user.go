package models

import "time"

// Auth represents authentication information including token and expiration time
type Auth struct {
	Token     string
	ExpiresIn time.Time
}

// User represents a user with credentials and authentication details
type User struct {
	Username       string `gorm:"primaryKey"`
	HashedPassword string
	Auth
}

// Credentials represents user login credentials
type Credentials struct {
	Username string
	Password string
}
