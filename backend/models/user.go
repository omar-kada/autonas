package models

import "time"

type Auth struct {
	Token     string
	ExpiresIn time.Time
}

type User struct {
	Username      string `gorm:"primaryKey"`
	HasedPassword string
	Auth
}

type Credentials struct {
	Username string
	Password string
}
