package balancer

import (
	"fmt"
	"net/http/httputil"
	"net/url"
	"sync"
)

type Backend struct {
	URL          *url.URL
	Alive        bool
	ReverseProxy *httputil.ReverseProxy
	mux          sync.RWMutex
}

func NewBackend(addr string) (*Backend, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %v", err)
	}

	return &Backend{
		URL:          u,
		Alive:        true,
		ReverseProxy: httputil.NewSingleHostReverseProxy(u),
	}, nil
}

func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.Alive = alive
}

func (b *Backend) IsAlive() bool {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.Alive
}
