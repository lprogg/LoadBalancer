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

	"github.com/lprogg/LoadBalancer/util"
)

var (
	port = flag.Int("port", util.Ports[0], "Starting port")
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
			fmt.Printf("Found service '%s' matching the specified request", service.Name)
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

	next := serviceList.NextServer()
	fmt.Printf("Forwarding to server: '%d'\n", next)
	
	serviceList.Servers[next].Proxy.ServeHTTP(w, r)
}

func InitNewLoadBalancer(c *util.Config) *LoadBalancer {
	listOfServers := make([]*util.Server, 0)
	mapOfServers := make(map[string]*util.ServerList, 0)

	for _, service := range c.Services {
		for _, replica := range service.Replicas {
			url, err := url.Parse(replica)

			if err != nil {
				log.Fatal(err)
			}

			proxy := httputil.NewSingleHostReverseProxy(url)
			
			listOfServers = append(listOfServers, &util.Server {
				Url: url,
				Proxy: proxy,
			})
		}
		mapOfServers[service.Matcher] = &util.ServerList{
			Servers: listOfServers,
			Current: 0,
			Name: service.Name,
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

	conf, err := util.LoadConfig(file)

	if err != nil {
		log.Fatal(err)
	}

	server := http.Server {
		Addr: fmt.Sprintf(":%d", *port),
		Handler: InitNewLoadBalancer(conf),
	}

	if err:= server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
