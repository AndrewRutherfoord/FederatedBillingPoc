package util

import (
	"context"
	"log"
	"time"

	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/repository"
	"github.com/andrewrutherfoord/fed-bill-poc/shared"
)

type Clock interface {
	Now() time.Time
}

type MockClock struct {
	currentTime  int64
	keyValueRepo repository.KeyValueRepository
	scheduler    *shared.Scheduler
}

func NewMockClock(startTime int64, keyValueRepo repository.KeyValueRepository, sched *shared.Scheduler) *MockClock {
	current, err := keyValueRepo.Get(context.Background(), "current_time")
	if err == nil {
		if parsed, err := time.Parse(time.RFC3339, current); err == nil {
			startTime = parsed.Unix()
		}
	}

	clock := &MockClock{currentTime: startTime, keyValueRepo: keyValueRepo, scheduler: sched}

	// Persist the initial time if it wasn't already set
	if err != nil {
		clock.persistCurrentTime()
	}

	log.Printf("MockClock initialized with time: %s", time.Unix(clock.currentTime, 0).UTC().Format(time.RFC3339))

	return clock
}

func (c *MockClock) persistCurrentTime() {
	c.keyValueRepo.Set(context.Background(), "current_time", time.Unix(c.currentTime, 0).UTC().Format(time.RFC3339))
}

func (c *MockClock) Now() time.Time {
	return time.Unix(c.currentTime, 0).UTC()
}

func (c *MockClock) Advance(seconds int64) {
	c.currentTime += seconds
	c.persistCurrentTime()

	// Execute scheduled jobs asynchronously so time advancement doesn't block
	go func() {
		if err := c.scheduler.CheckAndExecute(context.Background(), c.Now()); err != nil {
			log.Printf("error executing scheduled jobs: %v", err)
		}
	}()
}
