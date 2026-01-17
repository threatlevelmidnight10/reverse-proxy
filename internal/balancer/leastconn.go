package balancer

import (
	"math"
	"sync"

	"github.com/threatlevelmidnight10/reverse-proxy/internal/backend"
)

type LeastConn struct {
	backends []*backend.Backend
	mu       sync.RWMutex
}

func NewLeastConn() *LeastConn {
	return &LeastConn{
		backends: make([]*backend.Backend, 0),
	}
}

func (l *LeastConn) NextBackend() *backend.Backend {

	// actual least conn logic
	l.mu.RLock()
	defer l.mu.RUnlock()

	if len(l.backends) == 0 {
		return nil
	}

	var chosen *backend.Backend

	minConns := math.MaxInt64 - 1

	for _, b := range l.backends {
		if b.IsAlive() && (b.GetActiveConns() < int64(minConns)) {
			minConns = int(b.GetActiveConns())
			chosen = b
		}
	}

	return chosen
}

func (l *LeastConn) AddBackend(b *backend.Backend) {
	l.mu.Lock()

	defer l.mu.Unlock()

	l.backends = append(l.backends, b)
}

func (l *LeastConn) RemoveBackend(url string) {

	l.mu.Lock()
	defer l.mu.Unlock()

	for i, b := range l.backends {
		if b.URL.String() == url {
			l.backends = append(l.backends[:i], l.backends[i+1:]...)
			return
		}
	}
}
