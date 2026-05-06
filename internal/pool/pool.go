package pool

import (
	"fmt"
	"net/url"
	"sync"
)

type Backend struct {
	URL         *url.URL
	Alive       bool
	ActiveConns int64 // should be atomic.int64?
}

type Pool struct {
	Backends []*Backend
	mu       sync.RWMutex
}

func NewPool(rawURLs []string) (*Pool, error) {
	backends := make([]*Backend, 0, len(rawURLs))

	for _, raw := range rawURLs {
		url_, err := url.Parse(raw)
		if err != nil {
			return nil, fmt.Errorf("parse %q: %w", raw, err)
		}

		backends = append(backends, &Backend{URL: url_, Alive: true})
	}

	return &Pool{Backends: backends}, nil
}

func (p *Pool) Alive() []*Backend {
	p.mu.RLock()
	defer p.mu.RUnlock()

	alive := make([]*Backend, 0, len(p.Backends))
	for _, backnd := range p.Backends {
		if backnd.Alive {
			alive = append(alive, backnd)
		}
	}
	return alive
}

func (p *Pool) SetAlive(b *Backend, alive bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	b.Alive = alive
}
