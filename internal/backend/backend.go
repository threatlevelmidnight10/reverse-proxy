package backend

import (
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

// Backend server object
type Backend struct {
	URL         *url.URL // check?
	alive       bool
	ActiveConns int64        // Check?
	mu          sync.RWMutex //check?
	LastCheck   time.Time
	Weight      int
}

// create a new backend instance
func NewBackend(URL string, weight int) (*Backend, error) {
	urlStr, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}
	return &Backend{
		URL:    urlStr,
		Weight: weight,
		alive:  true,
	}, nil
}

func (b *Backend) isalive() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.alive
}

func (b *Backend) Setalive(alive bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.alive = alive
	b.LastCheck = time.Now()
}

func (b *Backend) IncrementConns() {
	atomic.AddInt64(&b.ActiveConns, 1)
}

func (b *Backend) DecrementConns() {
	atomic.AddInt64(&b.ActiveConns, -1)
}

func (b *Backend) GetActiveConns() int64 {
	return b.ActiveConns
}

// // Interface LB
// type LoadBalancer interface {

// 	// Get next backend host to connect to
// 	NextBackend() *Backend

// 	// Mark the health status of the backend it had connected to
// 	MarkBackendHealth(backendURL string, alive bool)
// }
