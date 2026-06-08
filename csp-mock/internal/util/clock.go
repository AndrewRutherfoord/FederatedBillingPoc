package util

import (
	"context"
	"log"
	"time"

	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/repository"
)

type Clock interface {
	Now() time.Time
}

type MockClock struct {
	currentTime  int64
	keyValueRepo repository.KeyValueRepository
}

func NewMockClock(startTime int64, keyValueRepo repository.KeyValueRepository) *MockClock {
	current, err := keyValueRepo.Get(context.Background(), "current_time")
	if err == nil {
		if parsed, err := time.Parse(time.RFC3339, current); err == nil {
			startTime = parsed.Unix()
		}
	}

	clock := &MockClock{currentTime: startTime, keyValueRepo: keyValueRepo}

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
}
