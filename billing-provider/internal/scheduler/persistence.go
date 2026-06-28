package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// KVSchedulerPersistence is an in-memory store of job last-execution times.
type KVSchedulerPersistence struct {
	mu    sync.RWMutex
	store map[string]string
}

func NewSchedulerPersistence() *KVSchedulerPersistence {
	return &KVSchedulerPersistence{
		store: make(map[string]string),
	}
}

func jobStateKey(jobID string) string {
	return fmt.Sprintf("scheduler:job:%s:last_execution", jobID)
}

func (p *KVSchedulerPersistence) LoadLastExecution(ctx context.Context, jobID string) (time.Time, error) {
	p.mu.RLock()
	val := p.store[jobStateKey(jobID)]
	p.mu.RUnlock()

	if val == "" {
		return time.Time{}, nil
	}

	return time.Parse(time.RFC3339, val)
}

func (p *KVSchedulerPersistence) SaveExecution(ctx context.Context, jobID string, t time.Time) error {
	p.mu.Lock()
	p.store[jobStateKey(jobID)] = t.UTC().Format(time.RFC3339)
	p.mu.Unlock()
	return nil
}
