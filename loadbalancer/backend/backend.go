package backend

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type Backend struct {
	URL          *url.URL
	Alive        bool
	Weight       int
	ResponseTime time.Duration
	mu           sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}

func NewBackend(rawURL string, weight int) (*Backend, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	
	backend := &Backend{
		URL:    parsedURL,
		Alive:  true,
		Weight: weight,
		ReverseProxy: httputil.NewSingleHostReverseProxy(parsedURL),
	}
	
	// Add error handling to reverse proxy
	backend.ReverseProxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		backend.SetAlive(false)
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "Unable to connect to backend: %v", err)
	}
	
	return backend, nil
}

func (b *Backend) SetAlive(alive bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Alive = alive
}

func (b *Backend) IsAlive() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.Alive
}

func (b *Backend) SetResponseTime(responseTime time.Duration) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.ResponseTime = responseTime
}

func (b *Backend) GetResponseTime() time.Duration {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.ResponseTime
}