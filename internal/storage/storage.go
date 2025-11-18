// Package storage defines all data store operations
package storage

// Storage is an abstraction of all database operations
type Storage interface {
	GetCurrentStacks() []string
}
