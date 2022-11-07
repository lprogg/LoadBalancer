package tests

import (
	"strings"
	"testing"

	"github.com/lprogg/LoadBalancer/util"
)

func TestLoadConfigWithRoundRobin(t *testing.T) {
	conf, err := util.LoadConfig(strings.NewReader(`
services:
  - 
    name: Round Robin Service
    matcher: /
    strategy: RoundRobin
    replicas:
      - url: http://localhost:8081
      - url: http://localhost:8082
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

	if conf.Services[0].Replicas[0].URL != "http://localhost:8081" {
		t.Errorf("First replica is expected to be 'http://localhost:8081', got instead: '%s'", conf.Services[0].Replicas[0])
	}

	if conf.Services[0].Replicas[1].URL != "http://localhost:8082" {
		t.Errorf("Second replica is expected to be 'http://localhost:8082', got instead: '%s'", conf.Services[0].Replicas[1])
	}
}

func TestLoadConfigWithWeightedRoundRobin(t *testing.T) {
	conf, err := util.LoadConfig(strings.NewReader(`
services:
  - 
    name: Weighted Round Robin Service
    matcher: /
    strategy: WeightedRoundRobin
    replicas:
      - url: http://localhost:8081
        metadata:
          weight: 10
      - url: http://localhost:8082
        metadata:
          weight: 5
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

	if conf.Services[0].Strategy != "WeightedRoundRobin" {
		t.Errorf("Strategy is expected to be 'WeightedRoundRobin', got instead: '%s'", conf.Services[0].Strategy)
	}

	if conf.Services[0].Name != "Weighted Round Robin Service" {
		t.Errorf("Service name is expected to be 'Weighted Round Robin Service', got instead: '%s'", conf.Services[0].Name)
	}

	if len(conf.Services[0].Replicas) != 2 {
		t.Errorf("Replicas are expected to be 2, got instead: '%d'", len(conf.Services[0].Replicas))
	}

	if conf.Services[0].Replicas[0].URL != "http://localhost:8081" {
		t.Errorf("First replica is expected to be 'localhost:8081', got instead: '%s'", conf.Services[0].Replicas[0])
	}

	if conf.Services[0].Replicas[1].URL != "http://localhost:8082" {
		t.Errorf("Second replica is expected to be 'localhost:8082', got instead: '%s'", conf.Services[0].Replicas[1])
	}

	if conf.Services[0].Replicas[0].Metadata["weight"] != "10" {
		t.Errorf("First replica is expected to have a weight of '10', got instead: '%s'", conf.Services[0].Replicas[0].Metadata["weight"])
	}

	if conf.Services[0].Replicas[1].Metadata["weight"] != "5" {
		t.Errorf("Second replica is expected to have a weight of '5', got instead: '%s'", conf.Services[0].Replicas[1].Metadata["weight"])
	}
}
