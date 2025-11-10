package backend

import (
	"sync"
)

type Pool struct {
	backends []*Backend
	mutex    sync.RWMutex
	current  int
}

func NewPool() *Pool {
	return &Pool{
		backends: make([]*Backend, 0),
		current:  0,
	}
}

func (p *Pool) AddBackend(rawURL string, weight int) error {
	backend, err := NewBackend(rawURL, weight)
	if err != nil {
		return err
	}
	
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.backends = append(p.backends, backend)
	return nil
}

func (p *Pool) GetBackends() []*Backend {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.backends
}

func (p *Pool) GetAliveBackends() []*Backend {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	
	alive := make([]*Backend, 0)
	for _, backend := range p.backends {
		if backend.IsAlive() {
			alive = append(alive, backend)
		}
	}
	return alive
}

func (p *Pool) MarkBackendStatus(backendURL string, alive bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	
	for _, b := range p.backends {
		if b.URL.String() == backendURL {
			b.SetAlive(alive)
			break
		}
	}
}