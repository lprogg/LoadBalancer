package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/lprogg/LoadBalancer/domain"
	"github.com/lprogg/LoadBalancer/ports"
	"github.com/lprogg/LoadBalancer/strategy"
	"github.com/lprogg/LoadBalancer/util"
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
	fmt.Printf("Trying to find matcher for the specified request '%s'\n", requestPath)
	for matcher, service := range lb.ServerList {
		if strings.HasPrefix(requestPath, matcher) {
			fmt.Printf("Found service '%s' matching the specified request\n", service.Name)
			return service, nil
		}
	}

	return nil, fmt.Errorf("could not find a matcher for the url: '%s'", requestPath)
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Received new request from url: '%s'\n", r.Host)
	
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

	fmt.Printf("Forwarding to server: '%s'\n\n", next.Url.Host)
	next.Proxy.ServeHTTP(w, r)
}

func InitNewLoadBalancer(c *util.Config) *LoadBalancer {
	listOfServers := make([]*domain.Server, 0)
	mapOfServers := make(map[string]*util.ServerList)

	for _, service := range c.Services {
		for _, replica := range service.Replicas {
			url, err := url.Parse(replica)

			if err != nil {
				log.Fatal(err)
			}

			proxy := httputil.NewSingleHostReverseProxy(url)
			
			listOfServers = append(listOfServers, &domain.Server {
				Url: url,
				Proxy: proxy,
			})
		}
		mapOfServers[service.Matcher] = &util.ServerList{
			Servers: listOfServers,
			Name: service.Name,
			Strategy: strategy.LoadStrategy(service.Strategy),
		}
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
