package modelsdb

import (
	"log/slog"
	"time"
)

//go:generate go run github.com/objectbox/objectbox-go/cmd/objectbox-gogen

// DeploymentStatus defines model for Deployment.Status.
type DeploymentStatus string

// Defines values for DeploymentStatus.
const (
	DeploymentStatusError   DeploymentStatus = "error"
	DeploymentStatusPlanned DeploymentStatus = "planned"
	DeploymentStatusRunning DeploymentStatus = "running"
	DeploymentStatusSuccess DeploymentStatus = "success"
)

// Deployment defines a deployment
type Deployment struct {
	Author  string
	Diff    string
	Events  []*Event    `objectbox:"link"`
	Files   []*FileDiff `objectbox:"link"`
	ID      uint64      `objectbox:"id"`
	Status  DeploymentStatus
	Time    time.Time `objectbox:"date"`
	EndTime time.Time `objectbox:"date"`
	Title   string
}

// FileDiff defines model for FileDiff.
type FileDiff struct {
	ID      uint64 `objectbox:"id"`
	Diff    string
	NewFile string
	OldFile string
}

// Event represent an event inside the deployment process
type Event struct {
	ID       uint64 `objectbox:"id"`
	Level    slog.Level
	Msg      string
	Time     time.Time `objectbox:"date"`
	ObjectID uint64
}
