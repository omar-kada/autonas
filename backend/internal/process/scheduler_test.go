package process

import (
	"omar-kada/autonas/internal/storage"
	"omar-kada/autonas/models"
	"path/filepath"
	"sync"
	"testing"
	"time"

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
	c := scheduler.Schedule(func(_ models.Config) {
		fnCalled <- true
	})

	// Verify the cron was created and started
	assert.NotNil(t, c, "Expected non-nil cron")
	assert.Equal(t, testConfig.CronPeriod, scheduler.currentPeriod, "Cron period mismatch")

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

func TestConcurrentScheduling(t *testing.T) {
	configStore := createTestConfigStore(t)
	scheduler := NewConfigScheduler(configStore).(*AtomicConfigScheduler)

	// Set up test config
	testConfig := models.Config{CronPeriod: "@every 1s"}
	err := configStore.Update(testConfig)
	assert.NoError(t, err)

	// Create a wait group to synchronize goroutines
	var wg sync.WaitGroup

	// Number of concurrent schedules
	numSchedules := 5

	// Create a channel to signal when all functions are called
	fnCalled := make(chan bool, numSchedules)

	// Schedule the function concurrently
	for range numSchedules {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c := scheduler.Schedule(func(_ models.Config) {
				fnCalled <- true
			})
			defer c.Stop()
		}()
	}

	// Wait for all schedules to complete
	wg.Wait()

	// Verify the number of function calls
	assert.Equal(t, numSchedules, len(fnCalled), "Number of function calls mismatch")
}

func TestGenerateAndRun(t *testing.T) {
	configStore := createTestConfigStore(t)
	scheduler := NewConfigScheduler(configStore).(*AtomicConfigScheduler)

	// Set up test config
	testConfig := models.Config{CronPeriod: "@every 1s", Services: map[string]models.ServiceConfig{
		"svc1": {},
	}}
	err := configStore.Update(testConfig)
	assert.NoError(t, err)

	// Create a channel to signal when the function is called
	fnCalled := make(chan bool, 1)

	// Call generateAndRun
	cfg := scheduler.generateAndRun(func(_ models.Config) {
		fnCalled <- true
	})

	// Verify the returned config
	assert.Equal(t, testConfig, cfg, "Config mismatch")

	// Verify the function was called
	select {
	case <-fnCalled:
		// Function was called
	default:
		t.Error("Function was not called")
	}
}
