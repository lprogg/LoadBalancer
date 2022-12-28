package health

import (
	"fmt"
	"net"
	"time"

	"github.com/lprogg/LoadBalancer/domain"
	log "github.com/sirupsen/logrus"
)

type HealthChecker struct {
	servers []*domain.Server
}

func InitNewHealthChecker(servers []*domain.Server) (*HealthChecker, error) {
	if len(servers) == 0 {
		return nil, fmt.Errorf("A server list expected, gotten an empty list")
	}

	return &HealthChecker{
		servers: servers,
	}, nil
}

func (hc *HealthChecker) Start() {
	log.Info("Starting health checker...\n")
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case _ = <- ticker.C:
			for _, server := range hc.servers {
				go checkHealth(server)
			}
		}
	}
}

func checkHealth(server *domain.Server) {
	_, err := net.DialTimeout("tcp", server.URL.Host, time.Second * 5)

	if err != nil {
		log.Errorf("Could not connect to the server at '%s'\n", server.URL.Host)
		
		initialState := server.SetServerLiveness(false)
		
		if initialState {
			log.Warnf("Set server '%s' from Live to Unavailable state\n", server.URL.Host)
		}

		return
	}

	initialState := server.SetServerLiveness(true)

	if !initialState {
		log.Infof("Set server '%s' from Unavailable to Live state\n", server.URL.Host)
	}
}
