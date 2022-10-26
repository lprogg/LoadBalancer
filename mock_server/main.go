package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/lprogg/LoadBalancer/ports"
)

var port = flag.Int("port", ports.MockServerPort, "Starting port of the mock server")

type MockServer struct {}

func (ms *MockServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("Server %d is ready.\n", *port)))
}

func main() {
	flag.Parse()
	
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), &MockServer{}); err != nil {
		log.Fatal(err)
	}
}
