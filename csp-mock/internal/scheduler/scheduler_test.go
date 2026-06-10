package scheduler

import (
	"context"
	"testing"
	"time"
)

// mockKVRepo is a simple in-memory key-value store for testing.
type mockKVRepo struct {
	data map[string]string
}

func (m *mockKVRepo) Get(ctx context.Context, key string) (string, error) {
	val, ok := m.data[key]
	if !ok {
		return "", nil
	}
	return val, nil
}

func (m *mockKVRepo) Set(ctx context.Context, key, value string) error {
	m.data[key] = value
	return nil
}

func (m *mockKVRepo) Delete(ctx context.Context, key string) error {
	delete(m.data, key)
	return nil
}

func TestCronSchedulerWithMockTime(t *testing.T) {
	ctx := context.Background()
	scheduler := New()

	executedTimes := []time.Time{}
	job := NewSimpleJob("hourly-batch", func(ctx context.Context) error {
		executedTimes = append(executedTimes, time.Now())
		return nil
	})

	// Register a job that runs at the top of every hour (0:MM:00)
	sched, err := NewCronSchedule("0 0 * * * *")
	if err != nil {
		t.Fatalf("failed to create cron schedule: %v", err)
	}

	if err := scheduler.Register(job, sched); err != nil {
		t.Fatalf("failed to register job: %v", err)
	}

	// Start at 12:00 on Jan 1, 2026
	start := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)

	// Initial check - should execute since no lastExecution time set
	if err := scheduler.CheckAndExecute(ctx, start); err != nil {
		t.Fatalf("failed to execute: %v", err)
	}

	if len(executedTimes) != 1 {
		t.Errorf("expected 1 execution, got %d", len(executedTimes))
	}

	// Advance 30 minutes - should NOT execute (next run is at 13:00)
	if err := scheduler.CheckAndExecute(ctx, start.Add(30*time.Minute)); err != nil {
		t.Fatalf("failed to execute: %v", err)
	}

	if len(executedTimes) != 1 {
		t.Errorf("expected 1 execution after 30 min advance, got %d", len(executedTimes))
	}

	// Advance to 13:00 - should execute
	if err := scheduler.CheckAndExecute(ctx, start.Add(1*time.Hour)); err != nil {
		t.Fatalf("failed to execute: %v", err)
	}

	if len(executedTimes) != 2 {
		t.Errorf("expected 2 executions after 1 hour advance, got %d", len(executedTimes))
	}

	// Call again at same time - should NOT re-execute
	if err := scheduler.CheckAndExecute(ctx, start.Add(1*time.Hour)); err != nil {
		t.Fatalf("failed to execute: %v", err)
	}

	if len(executedTimes) != 2 {
		t.Errorf("expected 2 executions after re-check at same time, got %d", len(executedTimes))
	}

	// Advance 2 hours - should execute once more
	if err := scheduler.CheckAndExecute(ctx, start.Add(3*time.Hour)); err != nil {
		t.Fatalf("failed to execute: %v", err)
	}

	if len(executedTimes) != 3 {
		t.Errorf("expected 3 executions after 3 hour advance, got %d", len(executedTimes))
	}
}

func TestSchedulerJobStatus(t *testing.T) {
	scheduler := New()
	job := NewSimpleJob("test-job", func(ctx context.Context) error { return nil })
	sched, _ := NewCronSchedule("0 0 * * * *")

	scheduler.Register(job, sched)

	lastExec, nextRun, ok := scheduler.GetJobStatus("test-job")
	if !ok {
		t.Fatal("job not found")
	}

	if !lastExec.IsZero() {
		t.Errorf("expected zero lastExecution, got %v", lastExec)
	}

	if nextRun.IsZero() {
		t.Errorf("expected non-zero nextRun, got %v", nextRun)
	}
}

func TestSchedulerPersistence(t *testing.T) {
	ctx := context.Background()
	kvRepo := &mockKVRepo{data: make(map[string]string)}

	// Create first scheduler, register and execute a job
	sched1 := NewWithPersistence(kvRepo)
	job := NewSimpleJob("persistent-job", func(ctx context.Context) error { return nil })
	cronSched, _ := NewCronSchedule("0 0 0 * * *") // Daily at midnight

	sched1.Register(job, cronSched)

	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	sched1.CheckAndExecute(ctx, start)

	// Verify the execution time was saved
	lastExec, nextRun, ok := sched1.GetJobStatus("persistent-job")
	if !ok {
		t.Fatal("job not found after execution")
	}
	if lastExec != start {
		t.Errorf("expected lastExec to be %v, got %v", start, lastExec)
	}

	// Create a second scheduler instance with the same kvRepo (simulating a restart)
	sched2 := NewWithPersistence(kvRepo)
	job2 := NewSimpleJob("persistent-job", func(ctx context.Context) error { return nil })
	sched2.Register(job2, cronSched)

	// The second scheduler should have loaded the execution time from the DB
	lastExec2, nextRun2, ok := sched2.GetJobStatus("persistent-job")
	if !ok {
		t.Fatal("job not found in reloaded scheduler")
	}

	if lastExec2 != start {
		t.Errorf("expected reloaded lastExec to be %v, got %v", start, lastExec2)
	}

	if nextRun != nextRun2 {
		t.Errorf("next run times should match: %v vs %v", nextRun, nextRun2)
	}

	// Check that executing again at the same time doesn't run (because we already ran at start)
	sched2.CheckAndExecute(ctx, start)
	lastExec2Again, _, _ := sched2.GetJobStatus("persistent-job")
	if lastExec2Again != start {
		t.Errorf("expected lastExec to remain %v, got %v", start, lastExec2Again)
	}

	// But advancing time should trigger a new execution
	nextDay := start.Add(24 * time.Hour)
	runCount := 0
	job3 := NewSimpleJob("persistent-job", func(ctx context.Context) error {
		runCount++
		return nil
	})

	sched3 := NewWithPersistence(kvRepo)
	sched3.Register(job3, cronSched)
	sched3.CheckAndExecute(ctx, nextDay)

	if runCount != 1 {
		t.Errorf("expected job to run once after time advance, got %d", runCount)
	}
}

func TestSchedulerMultipleJobs(t *testing.T) {
	ctx := context.Background()
	scheduler := New()

	job1Runs := 0
	job2Runs := 0

	job1 := NewSimpleJob("job-1", func(ctx context.Context) error {
		job1Runs++
		return nil
	})

	job2 := NewSimpleJob("job-2", func(ctx context.Context) error {
		job2Runs++
		return nil
	})

	sched1, _ := NewCronSchedule("0 0 * * * *") // Every hour at :00 seconds
	sched2, _ := NewCronSchedule("0 0 0 * * *") // Every day at 00:00:00

	scheduler.Register(job1, sched1)
	scheduler.Register(job2, sched2)

	start := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)

	// Both should run at start (both schedules can run at this time)
	scheduler.CheckAndExecute(ctx, start)
	if job1Runs != 1 || job2Runs != 1 {
		t.Errorf("expected both to run once, got job1=%d job2=%d", job1Runs, job2Runs)
	}

	// After 1 hour, only job1 should run
	scheduler.CheckAndExecute(ctx, start.Add(1*time.Hour))
	if job1Runs != 2 || job2Runs != 1 {
		t.Errorf("after 1h: expected job1=2 job2=1, got job1=%d job2=%d", job1Runs, job2Runs)
	}

	// After 24 hours, both should run again
	scheduler.CheckAndExecute(ctx, start.Add(24*time.Hour))
	if job1Runs != 3 && job2Runs != 2 {
		t.Errorf("after 24h: expected job1=3 job2=2, got job1=%d job2=%d", job1Runs, job2Runs)
	}
}
