package shared

import (
	"bufio"
	"log"
	"net/http"
	"sync"
	"time"
)

type Clock interface {
	Now() time.Time
}

type OnTimeAdvanceCallback interface {
	OnTimeAdvance(newTime time.Time)
}

type MockClock struct {
	currentTime time.Time
	mu          sync.RWMutex
	timeServer  string
	done        chan struct{}
	callbacks   []OnTimeAdvanceCallback
}

func NewMockClock(timeServerURL string, callbacks []OnTimeAdvanceCallback) *MockClock {
	clock := &MockClock{
		currentTime: time.Now(),
		timeServer:  timeServerURL,
		done:        make(chan struct{}),
		callbacks:   callbacks,
	}

	// Connect to time server SSE stream
	go clock.connectToTimeServer()

	log.Printf("MockClock initialized, connecting to time server at %s. Current time: %s", timeServerURL, clock.currentTime.UTC().Format(time.RFC3339))
	return clock
}

func (c *MockClock) connectToTimeServer() {
	for {
		resp, err := http.Get(c.timeServer + "/time/stream")
		if err != nil {
			log.Printf("Failed to connect to time server: %v, retrying in 5s", err)
			time.Sleep(5 * time.Second)
			continue
		}
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()

			// Parse SSE format: "data: <timestamp>"
			if len(line) < 6 || line[:5] != "data:" {
				continue
			}

			timeStr := line[6:]
			newTime, err := time.Parse(time.RFC3339Nano, timeStr)
			if err != nil {
				log.Printf("Failed to parse time from server: %v", err)
				continue
			}

			// Update local time
			c.mu.Lock()
			oldTime := c.currentTime
			c.currentTime = newTime
			c.mu.Unlock()

			log.Printf("Received clock update: %s", newTime.UTC().Format(time.RFC3339))

			// TODO: Implement a mechanism to trigger scheduled jobs
			// // If time actually advanced, trigger scheduled jobs
			if newTime.After(oldTime) {
				log.Printf("Time advanced to: %s", newTime.UTC().Format(time.RFC3339))
				for _, cb := range c.callbacks {
					go func() {
						cb.OnTimeAdvance(newTime)
					}()
				}
			}
		}

		if err := scanner.Err(); err != nil {
			log.Printf("SSE stream error: %v, reconnecting in 5s", err)
		}

		select {
		case <-c.done:
			return
		case <-time.After(5 * time.Second):
		}
	}
}

func (c *MockClock) Now() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.currentTime
}

// Shutdown gracefully closes the connection
func (c *MockClock) Shutdown() {
	close(c.done)
}
