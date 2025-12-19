package models

import (
	"log/slog"
	"time"
)

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
	ID      uint64 `gorm:"primaryKey;autoIncrement:true"`
	Author  string
	Diff    string
	Status  DeploymentStatus `gorm:"type:varchar(32)"`
	Time    time.Time        `gorm:"autoCreateTime"`
	EndTime time.Time
	Title   string
	Files   []FileDiff `gorm:"foreignKey:DeploymentID;constraint:OnDelete:CASCADE;"`
	Events  []Event    `gorm:"foreignKey:ObjectID;constraint:OnDelete:CASCADE;"`
}

// FileDiff defines model for FileDiff.
type FileDiff struct {
	ID           uint64 `gorm:"primaryKey;autoIncrement:true"`
	Diff         string
	NewFile      string
	OldFile      string
	DeploymentID uint64 `gorm:"index"`
}

// Event represent an event inside the deployment process
type Event struct {
	ID       uint64     `gorm:"primaryKey;autoIncrement:true"`
	Level    slog.Level `gorm:"type:int"`
	Msg      string
	Time     time.Time `gorm:"autoCreateTime"`
	ObjectID uint64    `gorm:"index"`
}

// Stats defines model for Stats.
type Stats struct {
	Author     string
	Error      int32
	LastDeploy time.Time
	LastStatus DeploymentStatus
	NextDeploy time.Time
	Success    int32
}
