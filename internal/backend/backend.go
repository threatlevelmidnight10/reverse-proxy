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
	activeConns int64        // Check?
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

func (b *Backend) IsAlive() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.alive
}

func (b *Backend) SetAlive(alive bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.alive = alive
	b.LastCheck = time.Now()
}

func (b *Backend) IncrementConns() {
	atomic.AddInt64(&b.activeConns, 1)
}

func (b *Backend) DecrementConns() {
	atomic.AddInt64(&b.activeConns, -1)
}

func (b *Backend) GetActiveConns() int64 {
	return atomic.LoadInt64(&b.activeConns)
}
