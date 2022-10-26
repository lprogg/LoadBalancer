package strategy

import (
	"sync/atomic"

	"github.com/lprogg/LoadBalancer/domain"
)

const (
	RoundRobinStrategy = "RoundRobin"
	WeightedRoundRobinStrategy = "WeightedRoundRobin"
	UnkownStrategy = "Unknown"
)

var strategies map[string]func() BalancingStrategy

type BalancingStrategy interface {
	NextServer([]*domain.Server) (*domain.Server, error)
}

type RoundRobin struct {
	current uint64
}

func (r *RoundRobin) NextServer(servers []*domain.Server) (*domain.Server, error) {
	next := atomic.AddUint64(&r.current, 1)
	serversLen := uint64(len(servers))
	return servers[next % serversLen], nil
}

func Init() {
	strategies = make(map[string]func() BalancingStrategy)
	strategies[RoundRobinStrategy] = func() BalancingStrategy {
		return &RoundRobin{current: 0}
	}
}

func LoadStrategy(name string) BalancingStrategy {
	strategy, ok := strategies[name]
	
	if !ok {
		return strategies[RoundRobinStrategy]()
	}

	return strategy()
}
