package storage

import (
	"omar-kada/autonas/modelsdb"
	"time"

	"github.com/objectbox/objectbox-go/objectbox"
)

// ObjectBoxStorage implements the Storage interface using ObjectBox
type ObjectBoxStorage struct {
	// ObjectBox store instance
	store *objectbox.ObjectBox
}

// NewObjectBoxStorage creates a new instance of ObjectBoxStorage
func NewObjectBoxStorage(store *objectbox.ObjectBox) *ObjectBoxStorage {
	return &ObjectBoxStorage{
		store: store,
	}
}

// DeploymentStorage implementation

// GetDeployments returns all deployments
func (s *ObjectBoxStorage) GetDeployments() ([]*modelsdb.Deployment, error) {
	box := modelsdb.BoxForDeployment(s.store)
	return box.GetAll()
}

// GetDeployment returns a specific deployment by ID
func (s *ObjectBoxStorage) GetDeployment(id uint64) (*modelsdb.Deployment, error) {
	box := modelsdb.BoxForDeployment(s.store)
	return box.Get(id)
}

// InitDeployment initializes a new deployment
func (s *ObjectBoxStorage) InitDeployment(title string, author string, diff string, files []*modelsdb.FileDiff) (*modelsdb.Deployment, error) {
	box := modelsdb.BoxForDeployment(s.store)
	filesBox := modelsdb.BoxForFileDiff(s.store)
	_, err := filesBox.PutMany(files)
	if err != nil {
		return nil, err
	}

	deployment := &modelsdb.Deployment{
		Title:  title,
		Author: author,
		Diff:   diff,
		Files:  files,
		Status: modelsdb.DeploymentStatusRunning,
		Time:   time.Now(),
		Events: []*modelsdb.Event{},
	}
	id, err := box.Put(deployment)
	if err != nil {
		return nil, err
	}
	deployment.ID = id
	return deployment, nil
}

// EndDeployment updates the status of a deployment
func (s *ObjectBoxStorage) EndDeployment(deploymentID uint64, status modelsdb.DeploymentStatus) error {
	box := modelsdb.BoxForDeployment(s.store)
	deployment, err := box.Get(deploymentID)
	if err != nil {
		return err
	}
	deployment.Status = status
	deployment.EndTime = time.Now()

	_, err = box.Put(deployment)
	return err
}

// EventStorage implementation

// StoreEvent stores a new event
func (s *ObjectBoxStorage) StoreEvent(event modelsdb.Event) error {
	eventBox := modelsdb.BoxForEvent(s.store)
	depBox := modelsdb.BoxForDeployment(s.store)
	_, err := eventBox.Put(&event)
	if err != nil {
		return err
	}
	dep, err := depBox.Get(event.ObjectID)
	if err != nil {
		return err
	}
	dep.Events = append(dep.Events, &event)
	_, err = depBox.Put(dep)
	return err
}

// GetEvents returns all events for a specific object ID
func (s *ObjectBoxStorage) GetEvents(objectID uint64) ([]*modelsdb.Event, error) {
	// Implementation for GetEvents
	box := modelsdb.BoxForEvent(s.store)
	return box.Query(modelsdb.Event_.ObjectID.Equals(objectID)).Find()
}
