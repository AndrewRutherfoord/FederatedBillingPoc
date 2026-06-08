package scheduler

import (
	"context"
	"log"
)

// SimpleJob is a basic job implementation for testing/examples.
type SimpleJob struct {
	id string
	fn func(context.Context) error
}

// NewSimpleJob creates a job that executes a function.
func NewSimpleJob(id string, fn func(context.Context) error) *SimpleJob {
	return &SimpleJob{id: id, fn: fn}
}

func (j *SimpleJob) ID() string {
	return j.id
}

func (j *SimpleJob) Execute(ctx context.Context) error {
	log.Printf("Executing job: %s", j.id)
	return j.fn(ctx)
}
