package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type TimeState struct {
	Current time.Time `json:"current"`
}

type TimeServer struct {
	current     time.Time
	mu          sync.RWMutex
	subscribers map[chan time.Time]struct{}
	subMu       sync.Mutex
	stateFile   string
}

func NewTimeServer(stateFile string) *TimeServer {
	ts := &TimeServer{
		current:     time.Now(),
		subscribers: make(map[chan time.Time]struct{}),
		stateFile:   stateFile,
	}
	ts.loadState()
	return ts
}

func (ts *TimeServer) loadState() {
	data, err := os.ReadFile(ts.stateFile)
	if err != nil {
		// File doesn't exist yet, that's fine
		log.Printf("No existing state file, starting fresh: %v", err)
		ts.saveState()
		return
	}

	var state TimeState
	if err := json.Unmarshal(data, &state); err != nil {
		log.Printf("Failed to load state: %v, starting fresh", err)
		ts.saveState()
		return
	}

	ts.current = state.Current
	log.Printf("Loaded time state: %v", ts.current)
}

func (ts *TimeServer) saveState() {
	ts.mu.RLock()
	state := TimeState{Current: ts.current}
	ts.mu.RUnlock()

	data, _ := json.MarshalIndent(state, "", "  ")
	if err := os.WriteFile(ts.stateFile, data, 0644); err != nil {
		log.Printf("Failed to save state: %v", err)
	}
}

func (ts *TimeServer) Subscribe() chan time.Time {
	ch := make(chan time.Time, 1)
	ts.subMu.Lock()
	ts.subscribers[ch] = struct{}{}
	ts.subMu.Unlock()

	ch <- ts.GetTime()
	return ch
}

func (ts *TimeServer) Unsubscribe(ch chan time.Time) {
	ts.subMu.Lock()
	delete(ts.subscribers, ch)
	ts.subMu.Unlock()
	close(ch)
}

func (ts *TimeServer) GetTime() time.Time {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	return ts.current
}

func (ts *TimeServer) Advance(duration time.Duration) {
	ts.mu.Lock()
	ts.current = ts.current.Add(duration)
	newTime := ts.current
	ts.mu.Unlock()

	// Persist to disk
	ts.saveState()

	// Broadcast to all subscribers
	ts.subMu.Lock()
	defer ts.subMu.Unlock()
	for ch := range ts.subscribers {
		select {
		case ch <- newTime:
		default:
		}
	}
}

func (ts *TimeServer) HandleSSE(w http.ResponseWriter, r *http.Request) {
	ch := ts.Subscribe()
	defer ts.Unsubscribe(ch)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher := w.(http.Flusher)

	for {
		select {
		case t := <-ch:
			fmt.Fprintf(w, "data: %s\n\n", t.Format(time.RFC3339Nano))
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}

func (ts *TimeServer) HandleAdvance(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Duration string `json:"duration"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	d, _ := time.ParseDuration(req.Duration)
	if d == 0 {
		d = 3600 * time.Second // Default to advancing 1 hour if no duration provided or parsing fails
	}

	ts.Advance(d)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]time.Time{"new_time": ts.GetTime()})
}

func main() {
	stateFile := "./time-state.json"
	ts := NewTimeServer(stateFile)

	http.HandleFunc("/time/stream", ts.HandleSSE)
	http.HandleFunc("/admin/advance", ts.HandleAdvance)

	log.Println("Mock time server running on :9999")
	http.ListenAndServe(":9999", nil)
}
