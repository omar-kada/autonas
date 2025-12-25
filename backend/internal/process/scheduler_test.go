package process

import (
	"path/filepath"
	"testing"
	"time"

	"omar-kada/autonas/internal/storage"
	"omar-kada/autonas/models"

	"github.com/stretchr/testify/assert"
)

// createTestConfigStore creates a ConfigStore with a temporary config file
func createTestConfigStore(t *testing.T) storage.ConfigStore {
	tempDir := t.TempDir()
	configFilePath := filepath.Join(tempDir, "config.yaml")
	return storage.NewConfigStore(configFilePath)
}

func TestNewConfigScheduler(t *testing.T) {
	configStore := createTestConfigStore(t)
	scheduler := NewConfigScheduler(configStore)

	assert.NotNil(t, scheduler, "Expected non-nil scheduler")
}

func TestSchedule(t *testing.T) {
	configStore := createTestConfigStore(t)
	scheduler := NewConfigScheduler(configStore).(*AtomicConfigScheduler)

	// Set up test config
	testConfig := models.Config{CronPeriod: "@every 1s"}
	err := configStore.Update(testConfig)
	assert.NoError(t, err)

	// Create a channel to signal when the function is called
	fnCalled := make(chan bool, 1)

	// Schedule the function
	c, err := scheduler.Schedule(func() {
		fnCalled <- true
	})

	// Verify the cron was created and started
	assert.NoError(t, err)
	assert.NotNil(t, c, "Expected non-nil cron")

	// Wait for the function to be called
	select {
	case <-fnCalled:
		// Function was called
	case <-time.After(2 * time.Second):
		t.Error("Function was not called within expected time")
	}

	// Stop the cron
	c.Stop()
}

func TestScheduleNoCronPeriod(t *testing.T) {
	configStore := createTestConfigStore(t)
	scheduler := NewConfigScheduler(configStore).(*AtomicConfigScheduler)

	// Set up test config with no cron period
	testConfig := models.Config{}
	err := configStore.Update(testConfig)
	assert.NoError(t, err)

	// Create a channel to signal when the function is called
	fnCalled := make(chan bool, 1)

	// Schedule the function
	c, err := scheduler.Schedule(func() {
		fnCalled <- true
	})

	// Verify the cron was not created
	assert.Error(t, err, "Expected error when no cron period is defined")
	assert.Nil(t, c, "Expected nil cron when no cron period is defined")

	// Verify the function is not called
	select {
	case <-fnCalled:
		t.Error("Function was called when no cron period was defined")
	case <-time.After(1 * time.Second):
		// Expected behavior - function not called
	}
}

func TestScheduleWhileRunning(t *testing.T) {
	configStore := createTestConfigStore(t)
	scheduler := NewConfigScheduler(configStore).(*AtomicConfigScheduler)

	// Set up test config
	testConfig := models.Config{CronPeriod: "@every 1s"}
	err := configStore.Update(testConfig)
	assert.NoError(t, err)

	// First function call channel
	fn1Called := make(chan bool, 1)
	// Second function call channel
	fn2Called := make(chan bool, 1)

	// Schedule first function
	c1, err := scheduler.Schedule(func() {
		fn1Called <- true
	})
	assert.NoError(t, err)
	assert.NotNil(t, c1, "Expected non-nil cron for first function")

	// Wait for first function to be called
	select {
	case <-fn1Called:
		// First function was called
	case <-time.After(2 * time.Second):
		t.Error("First function was not called within expected time")
	}

	// Schedule second function while first is running
	c2, err := scheduler.Schedule(func() {
		fn2Called <- true
	})
	assert.NoError(t, err)
	assert.NotNil(t, c2, "Expected non-nil cron for second function")

	// Verify both functions can be called
	select {
	case <-fn1Called:
		t.Error("First function was called again even when rescheduled")
	case <-time.After(2 * time.Second):
		t.Error("First function was not called again within expected time")
	case <-fn2Called:
	}

	// Stop both crons
	c2.Stop()
}
