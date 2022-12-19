package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/lprogg/LoadBalancer/domain"
	"github.com/lprogg/LoadBalancer/health"
	"github.com/lprogg/LoadBalancer/ports"
	"github.com/lprogg/LoadBalancer/strategy"
	"github.com/lprogg/LoadBalancer/util"
	log "github.com/sirupsen/logrus"
)

var (
	port = flag.Int("port", ports.DefaultPort, "Starting port")
	configFile = flag.String("config-path", "", "Config yaml file that needs to be supplied")
)

type LoadBalancer struct {
	Config *util.Config
	ServerList map[string]*util.ServerList
}

func (lb *LoadBalancer) findServiceList(requestPath string) (*util.ServerList, error) {
	log.Infof("Trying to find matcher for the specified request '%s'\n", requestPath)
	for matcher, service := range lb.ServerList {
		if strings.HasPrefix(requestPath, matcher) {
			log.Infof("Found service '%s' matching the specified request\n", service.Name)
			return service, nil
		}
	}

	return nil, fmt.Errorf("could not find a matcher for the url: '%s'", requestPath)
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Infof("Received new request from url: '%s'\n", r.Host)
	
	serviceList, err := lb.findServiceList(r.URL.Path)

	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	next, err := serviceList.Strategy.NextServer(serviceList.Servers)

	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Infof("Forwarding to server: '%s'\n", next.URL.Host)
	next.Proxy.ServeHTTP(w, r)
}

func InitNewLoadBalancer(c *util.Config) *LoadBalancer {
	listOfServers := make([]*domain.Server, 0)
	mapOfServers := make(map[string]*util.ServerList)

	for _, service := range c.Services {
		for _, replica := range service.Replicas {
			url, err := url.Parse(replica.URL)

			if err != nil {
				log.Fatal(err)
			}

			proxy := httputil.NewSingleHostReverseProxy(url)
			
			listOfServers = append(listOfServers, &domain.Server {
				URL: url,
				Proxy: proxy,
				Metadata: replica.Metadata,
			})
		}

		healthChecker, err := health.InitNewHealthChecker(listOfServers)

		if err != nil {
			log.Fatal(err)
		}

		mapOfServers[service.Matcher] = &util.ServerList{
			Servers: listOfServers,
			Name: service.Name,
			Strategy: strategy.LoadStrategy(service.Strategy),
			HealthChecker: healthChecker,
		}
	}

	for _, server := range mapOfServers {
		go server.HealthChecker.Start()
	}

	return &LoadBalancer {
		Config: c,
		ServerList: mapOfServers,
	}
}

func main() {
	flag.Parse()

	file, err := os.Open(*configFile)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	config, err := util.LoadConfig(file)

	if err != nil {
		log.Fatal(err)
	}

	server := http.Server {
		Addr: fmt.Sprintf(":%d", *port),
		Handler: InitNewLoadBalancer(config),
	}

	if err:= server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
