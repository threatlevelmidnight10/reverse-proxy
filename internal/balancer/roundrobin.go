package balancer

import (
	"sync"
	"sync/atomic"

	"github.com/threatlevelmidnight10/reverse-proxy/internal/backend"
)

type RoundRobin struct {
	counter  uint64
	backends []*backend.Backend
	mu       sync.RWMutex
}

func NewRoundRobin() *RoundRobin {
	return &RoundRobin{
		backends: make([]*backend.Backend, 0),
	}
}

func (r *RoundRobin) NextBackend() *backend.Backend {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// backendsAvailable
	var backendsAvailable []*backend.Backend

	for _, b := range r.backends {
		if b.IsAlive() {
			backendsAvailable = append(backendsAvailable, b)
		}
	}

	if len(backendsAvailable) == 0 {
		return nil
	}

	count := atomic.AddUint64(&r.counter, 1)
	nextBackendIndex := (count - 1) % uint64(len(backendsAvailable))

	return backendsAvailable[nextBackendIndex]
}

func (r *RoundRobin) AddBackend(b *backend.Backend) {
	r.mu.Lock()

	defer r.mu.Unlock()

	r.backends = append(r.backends, b)
}

func (r *RoundRobin) RemoveBackend(url string) {

	r.mu.Lock()
	defer r.mu.Unlock()

	for i, b := range r.backends {
		if b.URL.String() == url {
			r.backends = append(r.backends[:i], r.backends[i+1:]...)
			return
		}
	}
}
