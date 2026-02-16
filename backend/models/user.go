package models

import "time"

// Session represents a user session with authentication details
type Session struct {
	SessionID        uint64 `gorm:"primaryKey;autoGenerate"`
	RefreshToken     string `gorm:"index"`
	RefreshExpiresIn time.Time
	Revoked          bool
	Username         string `gorm:"not null"`
}

// User represents a user with credentials and authentication details
type User struct {
	Username       string `gorm:"primaryKey"`
	HashedPassword string
}

// Credentials represents user login credentials
type Credentials struct {
	Username string
	Password string
}
