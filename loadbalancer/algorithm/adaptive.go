package algorithm

import (
	"loadbalancer/backend"
	"sync/atomic"
)

type RoundRobin struct {
	pool    *backend.Pool
	current uint64
}

func NewRoundRobin(pool *backend.Pool) *RoundRobin {
	return &RoundRobin{
		pool:    pool,
		current: 0,
	}
}

func (r *RoundRobin) GetNextBackend() *backend.Backend {
	aliveBackends := r.pool.GetAliveBackends()
	if len(aliveBackends) == 0 {
		return nil
	}
	
	next := atomic.AddUint64(&r.current, 1)
	index := int(next) % len(aliveBackends)
	return aliveBackends[index]
}

func (r *RoundRobin) Name() string {
	return "round-robin"
}