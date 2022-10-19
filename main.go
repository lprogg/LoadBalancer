package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/lprogg/LoadBalancer/util"
)

var (
	port = flag.Int("port", 8080, "Starting port")
	configFile = flag.String("config-path", "", "Config yaml file that needs to be supplied")
)

type LoadBalancer struct {
	Config *util.Config
	ServerList *util.ServerList
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Received new request from url: '%s'\n", r.Host)
	
	next := lb.ServerList.NextServer()
	fmt.Printf("Forwarding to server: '%d'\n", next)
	
	lb.ServerList.Servers[next].Proxy.ServeHTTP(w, r)
}

func InitNewLoadBalancer(c *util.Config) *LoadBalancer {
	listOfServers := make([]*util.Server, 0)

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
	}

	return &LoadBalancer {
		Config: c,
		ServerList: &util.ServerList {
			Servers: listOfServers,
			Current: 0,
		},
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
