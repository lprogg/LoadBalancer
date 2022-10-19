package mock_server

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var port = flag.Int("port", 8081, "Port to start the demo server")

type DemoServer struct {}

func (ds *DemoServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("Server %d is ready.\n", *port)))
}

func main() {
	flag.Parse()
	
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), &DemoServer{}); err != nil {
		log.Fatal(err)
	}
}
