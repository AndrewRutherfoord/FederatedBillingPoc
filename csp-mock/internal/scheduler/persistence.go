package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/repository"
)

type KVSchedulerPersistence struct {
	kvRepo repository.KeyValueRepository
}

func NewSchedulerPersistence(kvRepo repository.KeyValueRepository) *KVSchedulerPersistence {
	return &KVSchedulerPersistence{kvRepo: kvRepo}
}

func jobStateKey(jobID string) string {
	return fmt.Sprintf("scheduler:job:%s:last_execution", jobID)
}

func (p *KVSchedulerPersistence) LoadLastExecution(ctx context.Context, jobID string) (time.Time, error) {
	val, err := p.kvRepo.Get(ctx, jobStateKey(jobID))
	if err != nil || val == "" {
		return time.Time{}, err
	}

	t, err := time.Parse(time.RFC3339, val)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func (p *KVSchedulerPersistence) SaveExecution(ctx context.Context, jobID string, t time.Time) error {
	return p.kvRepo.Set(ctx, jobStateKey(jobID), t.UTC().Format(time.RFC3339))
}
