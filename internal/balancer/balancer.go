package balancer

import (
	"github.com/threatlevelmidnight10/reverse-proxy/internal/backend"
)

// Interface LB
type LoadBalancer interface {

	// Get next backend host to send rqeuest to, returns nil if no backend available
	NextBackend() *backend.Backend

	AddBackend(b *backend.Backend)

	RemoveBackend(url string)
}
