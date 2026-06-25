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

// Render Basic Web Interface to advance clock
const indexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>Mock Time Server</title>
<style>
  body { font-family: sans-serif; text-align: center; margin-top: 4rem; }
  #clock { font-size: 3rem; margin-bottom: 2rem; }
  button {
    font-size: 1.1rem;
    padding: 0.6rem 1.2rem;
    margin: 0.3rem;
    cursor: pointer;
  }
</style>
</head>
<body>
  <div id="clock">--:--:--</div>
  <div>
    <button onclick="advance('5m')">+5 min</button>
    <button onclick="advance('15m')">+15 min</button>
    <button onclick="advance('30m')">+30 min</button>
    <button onclick="advance('60m')">+60 min</button>
    <button onclick="advance('90m')">+90 min</button>
  </div>

<script>
  const clock = document.getElementById('clock');

  function render(isoString) {
    const d = new Date(isoString);
    clock.textContent = d.toLocaleString();
  }

  function advance(duration) {
    fetch('/admin/advance', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ duration })
    })
      .then(res => res.json())
      .then(data => render(data.new_time))
      .catch(err => console.error('Failed to advance time', err));
  }

  function connect() {
    const es = new EventSource('/time/stream');
    es.onmessage = (e) => render(e.data);
    es.onerror = () => {
      es.close();
      setTimeout(connect, 2000);
    };
  }
  connect();
</script>
</body>
</html>
`

func (ts *TimeServer) HandleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(indexHTML))
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

	http.HandleFunc("/", ts.HandleIndex)
	http.HandleFunc("/time/stream", ts.HandleSSE)
	http.HandleFunc("/admin/advance", ts.HandleAdvance)

	log.Println("Mock time server running on :9999")
	http.ListenAndServe(":9999", nil)
}
