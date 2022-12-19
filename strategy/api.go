package strategy

import (
	"sync"
	"sync/atomic"

	"github.com/lprogg/LoadBalancer/domain"
	log "github.com/sirupsen/logrus"
)

const (
	RoundRobinStrategy = "RoundRobin"
	WeightedRoundRobinStrategy = "WeightedRoundRobin"
)

var strategies map[string]func() BalancingStrategy

type BalancingStrategy interface {
	NextServer([]*domain.Server) (*domain.Server, error)
}

type RoundRobin struct {
	current uint64
}

type WeightedRoundRobin struct {
	mutex sync.Mutex
	count []int
	current int
}

func (r *RoundRobin) NextServer(servers []*domain.Server) (*domain.Server, error) {
	next := atomic.AddUint64(&r.current, 1)
	serversLen := uint64(len(servers))
	selected := servers[next % serversLen]

	log.Infof("Strategy selected server: '%s'\n", selected.URL.Host)

	return selected, nil
}

func (r *WeightedRoundRobin) NextServer(servers []*domain.Server) (*domain.Server, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.count == nil {
		r.count = make([]int, len(servers))
		r.current = 0
	}

	capacity := servers[r.current].GetMetadataOrDefaultInt("weight", 1)
	
	if r.count[r.current] <= capacity {
		r.count[r.current] += 1
		log.Infof("Strategy selected server: '%s'\n", servers[r.current].URL.Host)
		return servers[r.current], nil
	}

	r.count[r.current] = 0
	r.current = (r.current + 1) % len(servers)
	log.Infof("Strategy selected server: '%s'\n", servers[r.current].URL.Host)
	
	return servers[r.current], nil
}

func init() {
	strategies = make(map[string]func() BalancingStrategy)
	
	strategies[RoundRobinStrategy] = func() BalancingStrategy {
		return &RoundRobin{current: 0}
	}

	strategies[WeightedRoundRobinStrategy] = func() BalancingStrategy {
		return &WeightedRoundRobin{mutex: sync.Mutex{}}
	}
}

func LoadStrategy(name string) BalancingStrategy {
	strategy, ok := strategies[name]
	
	if !ok {
		log.Warnf("Strategy '%s' not found, falling back to a RoundRobinStrategy\n\n", name)
		return strategies[RoundRobinStrategy]()
	}

	log.Infof("Picked strategy '%s'\n\n", name)
	return strategy()
}
