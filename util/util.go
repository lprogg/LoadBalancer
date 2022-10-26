package util

import (
	"github.com/lprogg/LoadBalancer/domain"
	"github.com/lprogg/LoadBalancer/strategy"
)

type Config struct {
	Services []*domain.Service `yml:"services"`
}

type ServerList struct {
	Servers []*domain.Server
	Name string
	Strategy strategy.BalancingStrategy
}
