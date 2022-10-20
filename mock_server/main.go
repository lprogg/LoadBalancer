package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/lprogg/LoadBalancer/util"
)

var port = flag.Int("port", util.Ports[1], "Starting port of the mock server")

type MockServer struct {}

func (ds *MockServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("Server %d is ready.\n", *port)))
}

func main() {
	flag.Parse()
	
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), &MockServer{}); err != nil {
		log.Fatal(err)
	}
}
