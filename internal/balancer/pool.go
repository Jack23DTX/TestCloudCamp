package balancer

import (
	"log"
	"net"
	"sync"
	"time"
)

type ServerPool struct {
	backends []*Backend
	current  int
	mu       sync.Mutex
}

func NewServerPool(backends []*Backend) *ServerPool {
	return &ServerPool{
		backends: backends,
	}
}

func (s *ServerPool) GetNextBackend() *Backend {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := 0; i < len(s.backends); i++ {
		next := (s.current + 1) % len(s.backends)
		s.current = next
		if s.backends[next].IsAlive() {
			return s.backends[next]
		}
	}
	return nil
}

func (s *ServerPool) StartHealthChecks(interval time.Duration) {
	ticker := time.NewTicker(interval)

	defer ticker.Stop()

	for range ticker.C {
		s.checkAllBackend()
	}
}

func (s *ServerPool) checkAllBackend() {
	s.mu.Lock()
	backends := s.backends
	s.mu.Unlock()

	for _, b := range backends {
		go s.checkBackend(b)
	}
}

func (s *ServerPool) checkBackend(b *Backend) {

	conn, err := net.DialTimeout("tcp", b.URL.Host, 2*time.Second)
	if err != nil {
		log.Printf("Health check FAILED for %s: %v", b.URL.Host, err)
		b.SetAlive(false)
		return
	}
	if err := conn.Close(); err != nil {
		return
	}
	log.Printf("Health check SUCCESS for %s", b.URL.Host)
	b.SetAlive(true)
}
