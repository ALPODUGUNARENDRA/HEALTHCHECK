package algorithm

import "loadbalancer/backend"

type LoadBalancerAlgorithm interface {
	GetNextBackend() *backend.Backend
	Name() string
}