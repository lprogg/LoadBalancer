package strategy

import (
	"fmt"
	"sync/atomic"

	"github.com/lprogg/LoadBalancer/domain"
)

const RoundRobinStrategy = "RoundRobin"

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
	selected := servers[next % serversLen]

	fmt.Printf("Strategy selected server: '%s'\n", selected.Url.Host)

	return selected, nil
}

func init() {
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
