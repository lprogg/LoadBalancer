package strategy

import (
	"fmt"
	"sync"

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
	mutex sync.Mutex
	current int
}

type WeightedRoundRobin struct {
	mutex sync.Mutex
	count []int
	current int
}

func (r *RoundRobin) NextServer(servers []*domain.Server) (*domain.Server, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	seen := 0
	var selected *domain.Server

	for seen < len(servers) {
		selected = servers[r.current]
		r.current = (r.current + 1) % len(servers)
		if selected.IsAlive() {
			break
		}

		seen += 1
	}

	if selected == nil || seen == len(servers) {
		log.Error("All servers are down")
		return nil, fmt.Errorf("Checked all the '%d' servers, none available", seen)
	}

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

	seen := 0
	var selected *domain.Server

	for seen < len(servers) {
		selected = servers[r.current]
		capacity := selected.GetMetadataOrDefaultInt("weight", 1)
		
		if !selected.IsAlive() {
			seen += 1
			r.count[r.current] = 0
			r.current = (r.current + 1) % len(servers)
			continue
		}
		
		if r.count[r.current] <= capacity {
			r.count[r.current] += 1
			log.Infof("Strategy selected server: '%s'\n", selected.URL.Host)
			return selected, nil
		}

		r.count[r.current] = 0
		r.current = (r.current + 1) % len(servers)
	}

	if selected == nil || seen == len(servers) {
		log.Error("All servers are down")
		return nil, fmt.Errorf("Checked all the '%d' servers, none available", seen)
	}

	return selected, nil
}

func init() {
	strategies = make(map[string]func() BalancingStrategy)
	
	strategies[RoundRobinStrategy] = func() BalancingStrategy {
		return &RoundRobin{
			mutex: sync.Mutex{},
			current: 0,
		}
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

	log.Infof("Selected strategy '%s'\n\n", name)
	return strategy()
}
