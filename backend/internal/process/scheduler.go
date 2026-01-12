package process

import (
	"fmt"
	"sync"
	"time"

	"omar-kada/autonas/internal/storage"

	"github.com/robfig/cron/v3"
)

// ConfigScheduler is responsible for cron running a sheduled job with updated config
type ConfigScheduler interface {
	Schedule(fn func()) (*cron.Cron, error)
	GetNext() time.Time
}

// NewConfigScheduler creates a new ConfigScheduler that ensures only one cron job runs at a time.
func NewConfigScheduler(configStore storage.ConfigStore) ConfigScheduler {
	return &AtomicConfigScheduler{
		configStore: configStore,
	}
}

// AtomicConfigScheduler runs only a single cron job at a time
type AtomicConfigScheduler struct {
	configStore storage.ConfigStore
	cron        *cron.Cron
	mu          sync.Mutex
}

// Schedule stops the old cron when it exists, and runs a new cron job
func (a *AtomicConfigScheduler) Schedule(fn func()) (*cron.Cron, error) {
	// make sure only one sync job is running at a time
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.cron != nil {
		a.cron.Stop()
		a.cron = nil
	}
	c := cron.New()

	cfg, err := a.configStore.Get()
	if err != nil {
		return nil, err
	}
	if cfg.CronPeriod != "" {

		_, err := c.AddFunc(cfg.CronPeriod, fn)
		if err != nil {
			return nil, err
		}
		c.Start()
		a.cron = c
		return c, nil
	}

	fn()
	return nil, fmt.Errorf("couldn't schedule job, no cron period is defined")
}

// GetNext returns the next scheduled time of the cron job.
// If no cron job is scheduled or no entries are present, it returns the zero time.
func (a *AtomicConfigScheduler) GetNext() time.Time {
	if a.cron == nil {
		return time.Time{}
	}
	entries := a.cron.Entries()
	if len(entries) == 0 {
		return time.Time{}
	}
	return entries[0].Next
}
