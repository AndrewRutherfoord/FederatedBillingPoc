package shared

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

// Job represents a scheduled job.
type Job interface {
	ID() string
	Execute(ctx context.Context, startTime time.Time) error
}

// Schedule describes when a job should run.
type Schedule interface {
	Next(t time.Time) time.Time
}

// CronSchedule wraps a cron expression.
type CronSchedule struct {
	sched cron.Schedule
}

// NewCronSchedule creates a new cron-based schedule from an expression.
// Cron format: "second minute hour day month weekday"
// Examples: "@hourly", "0 0 * * *" (daily at midnight), "0 */15 * * * *" (every 15 mins)
func NewCronSchedule(cronExpr string) (*CronSchedule, error) {
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	sched, err := parser.Parse(cronExpr)
	if err != nil {
		return nil, err
	}
	return &CronSchedule{sched: sched}, nil
}

func (cs *CronSchedule) Next(t time.Time) time.Time {
	return cs.sched.Next(t)
}

// jobEntry tracks a job and its schedule.
type jobEntry struct {
	job           Job
	schedule      Schedule
	lastExecution time.Time
}

// Persistence layer needs to be implemented by the caller, as it involves an external system
type SchedulerPersistence interface {
	LoadLastExecution(ctx context.Context, jobID string) (time.Time, error)
	SaveExecution(ctx context.Context, jobID string, t time.Time) error
}

// Scheduler manages scheduled jobs and executes them based on time.
type Scheduler struct {
	mu          sync.RWMutex
	jobs        map[string]*jobEntry
	persistence SchedulerPersistence
}

// New creates a new scheduler without persistence.
func New() *Scheduler {
	return &Scheduler{
		jobs: make(map[string]*jobEntry),
	}
}

// NewWithPersistence creates a scheduler that persists job execution times to a key-value store.
func NewWithPersistence(persistence SchedulerPersistence) *Scheduler {
	return &Scheduler{
		jobs:        make(map[string]*jobEntry),
		persistence: persistence,
	}
}

// loadLastExecution loads the last execution time for a job from the key-value store.
func (s *Scheduler) loadLastExecution(ctx context.Context, jobID string) time.Time {
	if s.persistence == nil {
		return time.Time{}
	}

	val, err := s.persistence.LoadLastExecution(ctx, jobID)
	if err != nil || val == (time.Time{}) {
		return time.Time{}
	}

	return val
}

// saveLastExecution saves the last execution time for a job to the key-value store.
func (s *Scheduler) saveLastExecution(ctx context.Context, jobID string, t time.Time) error {
	if s.persistence == nil {
		return nil
	}

	return s.persistence.SaveExecution(ctx, jobID, t)
}

// Register adds a job to the scheduler.
func (s *Scheduler) Register(job Job, schedule Schedule) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.jobs[job.ID()]; exists {
		return fmt.Errorf("job %q already registered", job.ID())
	}

	// Load the last execution time from persistence if available
	lastExec := s.loadLastExecution(context.Background(), job.ID())

	s.jobs[job.ID()] = &jobEntry{
		job:           job,
		schedule:      schedule,
		lastExecution: lastExec,
	}
	return nil
}

// Unregister removes a job from the scheduler.
func (s *Scheduler) Unregister(jobID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.jobs[jobID]; !exists {
		return fmt.Errorf("job %q not registered", jobID)
	}

	delete(s.jobs, jobID)
	return nil
}

// CheckAndExecute evaluates all jobs and executes those whose scheduled time has arrived.
// This is called after the mock clock advances, or on a timer in production.
func (s *Scheduler) CheckAndExecute(ctx context.Context, now time.Time) error {
	s.mu.RLock()
	entries := make([]*jobEntry, 0, len(s.jobs))
	for _, entry := range s.jobs {
		entries = append(entries, entry)
	}
	s.mu.RUnlock()

	for _, entry := range entries {
		nextRun := entry.schedule.Next(entry.lastExecution)
		// If the next scheduled time has arrived and we haven't executed yet, run the job.
		if nextRun.Before(now) || nextRun.Equal(now) {
			jobID := entry.job.ID()
			log.Printf("[scheduler] starting job: %s", jobID)
			start := time.Now()

			if err := entry.job.Execute(ctx, now); err != nil {
				return fmt.Errorf("job %q failed: %w", jobID, err)
			}

			duration := time.Since(start)
			log.Printf("[scheduler] completed job: %s (took %v)", jobID, duration)

			s.mu.Lock()
			entry.lastExecution = now
			s.mu.Unlock()

			// Persist the execution time so it survives restarts
			if err := s.saveLastExecution(ctx, jobID, now); err != nil {
				return fmt.Errorf("failed to persist job %q execution: %w", jobID, err)
			}
		}
	}
	return nil
}

// GetJobStatus returns info about a registered job (for debugging/testing).
func (s *Scheduler) GetJobStatus(jobID string) (lastExecution time.Time, nextRun time.Time, ok bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, exists := s.jobs[jobID]
	if !exists {
		return time.Time{}, time.Time{}, false
	}

	return entry.lastExecution, entry.schedule.Next(entry.lastExecution), true
}

func (s *Scheduler) OnTimeAdvance(newTime time.Time) {
	if err := s.CheckAndExecute(context.Background(), newTime); err != nil {
		log.Printf("error executing scheduled jobs: %v", err)
	}
}

type JobToRegister struct {
	job      Job
	cronExpr string
}

func NewJobToRegister(job Job, cronExpr string) JobToRegister {
	return JobToRegister{job: job, cronExpr: cronExpr}
}

func RegisterJobs(s *Scheduler, jobs []JobToRegister) error {
	for _, rj := range jobs {
		sched, err := NewCronSchedule(rj.cronExpr)
		if err != nil {
			return fmt.Errorf("invalid cron expression for job %q: %w", rj.job.ID(), err)
		}
		if err := s.Register(rj.job, sched); err != nil {
			return fmt.Errorf("failed to register job %q: %w", rj.job.ID(), err)
		}
	}
	return nil
}
