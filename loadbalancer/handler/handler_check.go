package handler

import (
	"loadbalancer/backend"
	"net/http"
	"time"
)

type HealthChecker struct {
	pool     *backend.Pool
	interval time.Duration
	stop     chan bool
}

func NewHealthChecker(pool *backend.Pool, interval time.Duration) *HealthChecker {
	return &HealthChecker{
		pool:     pool,
		interval: interval,
		stop:     make(chan bool),
	}
}

func (h *HealthChecker) Start() {
	ticker := time.NewTicker(h.interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				h.checkHealth()
			case <-h.stop:
				ticker.Stop()
				return
			}
		}
	}()
}

func (h *HealthChecker) Stop() {
	h.stop <- true
}

func (h *HealthChecker) checkHealth() {
	backends := h.pool.GetBackends()
	
	for _, backend := range backends {
		go func(b *backend.Backend) {
			start := time.Now()
			
			resp, err := http.Get(b.URL.String() + "/health")
			if err != nil || resp.StatusCode != http.StatusOK {
				b.SetAlive(false)
				return
			}
			defer resp.Body.Close()
			
			responseTime := time.Since(start)
			b.SetResponseTime(responseTime)
			b.SetAlive(true)
		}(backend)
	}
}