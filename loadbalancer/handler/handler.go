package handler

import (
	"loadbalancer/algorithm"
	"loadbalancer/backend"
	"net/http"
)

type LoadBalancerHandler struct {
	pool      *backend.Pool
	algorithm algorithm.LoadBalancerAlgorithm
}

type LoadBalancerAlgorithm interface {
	GetNextBackend() *backend.Backend
	Name() string
}

func NewLoadBalancerHandler(pool *backend.Pool, algorithmType string) *LoadBalancerHandler {
	var algo LoadBalancerAlgorithm
	
	switch algorithmType {
	case "adaptive":
		algo = algorithm.NewAdaptive(pool)
	default:
		algo = algorithm.NewRoundRobin(pool)
	}
	
	return &LoadBalancerHandler{
		pool:      pool,
		algorithm: algo,
	}
}

func (h *LoadBalancerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	backend := h.algorithm.GetNextBackend()
	
	if backend == nil {
		http.Error(w, "No available backend servers", http.StatusServiceUnavailable)
		return
	}
	
	// Add headers for debugging
	w.Header().Set("X-Load-Balancer", "Go-LB")
	w.Header().Set("X-Backend", backend.URL.Host)
	
	// Proxy the request
	backend.ReverseProxy.ServeHTTP(w, r)
}