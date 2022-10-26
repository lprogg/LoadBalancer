package tests

import (
	"strings"
	"testing"

	"github.com/lprogg/LoadBalancer/util"
)

func TestLoadConfig(t *testing.T) {
	conf, err := util.LoadConfig(strings.NewReader(`
services:
  - 
    name: Round Robin Service
    matcher: /
    strategy: RoundRobin
    replicas:
      - localhost:8081
      - localhost:8082
`))

	if err != nil {
		t.Errorf("Error should be nil: '%s'", err)
	}

	if len(conf.Services) != 1 {
		t.Errorf("Expected to be 1 service, got instead: '%d'", len(conf.Services))
	}

	if conf.Services[0].Matcher != "/" {
		t.Errorf("Matcher is expected to be '/', got instead: '%s'", conf.Services[0].Matcher)
	}

	if conf.Services[0].Strategy != "RoundRobin" {
		t.Errorf("Strategy is expected to be 'RoundRobin', got instead: '%s'", conf.Services[0].Strategy)
	}

	if conf.Services[0].Name != "Round Robin Service" {
		t.Errorf("Service name is expected to be 'Round Robin Service', got instead: '%s'", conf.Services[0].Name)
	}

	if len(conf.Services[0].Replicas) != 2 {
		t.Errorf("Replicas are expected to be 2, got instead: '%d'", len(conf.Services[0].Replicas))
	}

	if conf.Services[0].Replicas[0] != "localhost:8081" {
		t.Errorf("First replica is expected to be 'localhost:8081', got instead: '%s'", conf.Services[0].Replicas[0])
	}

	if conf.Services[0].Replicas[1] != "localhost:8082" {
		t.Errorf("Second replica is expected to be 'localhost:8082', got instead: '%s'", conf.Services[0].Replicas[1])
	}
}
