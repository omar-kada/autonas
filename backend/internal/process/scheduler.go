package process

import (
	"fmt"
	"log/slog"
	"omar-kada/autonas/internal/storage"
	"omar-kada/autonas/models"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

// ConfigScheduler is responsible for cron running a sheduled job with updated config
type ConfigScheduler interface {
	Schedule(fn func(cfg models.Config)) *cron.Cron
	getNext() time.Time
}

// NewConfigScheduler creates a new ConfigScheduler that ensures only one cron job runs at a time.
func NewConfigScheduler(configStore storage.ConfigStore) ConfigScheduler {
	return &AtomicConfigScheduler{
		configStore: configStore,
	}
}

// AtomicConfigScheduler runs only a single cron job at a time
type AtomicConfigScheduler struct {
	configStore   storage.ConfigStore
	currentPeriod string
	cron          *cron.Cron
	mu            sync.Mutex
}

// Schedule stops the old cron when it exists, and runs a new cron job
func (a *AtomicConfigScheduler) Schedule(fn func(cfg models.Config)) *cron.Cron {
	// make sure only one sync job is running at a time
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.cron != nil {
		a.cron.Stop()
		a.cron = nil
	}
	c := cron.New()

	cfg := a.generateAndRun(fn)
	a.currentPeriod = cfg.CronPeriod

	c.AddFunc(a.currentPeriod, func() {
		a.generateAndRun(fn)
	})

	c.Start()
	a.cron = c
	return c
}

func (a *AtomicConfigScheduler) generateAndRun(fn func(cfg models.Config)) models.Config {
	cfg, err := a.configStore.Get()
	if err != nil {
		slog.Error(fmt.Sprintf("couldn't get config : %v", err))
	}
	fn(cfg)
	return cfg
}

func (a *AtomicConfigScheduler) getNext() time.Time {
	if a.cron == nil {
		return time.Time{}
	}
	entries := a.cron.Entries()
	if len(entries) == 0 {
		return time.Time{}
	}
	return entries[0].Next
}
