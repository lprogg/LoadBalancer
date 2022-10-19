package util

import (
	"net/http/httputil"
	"net/url"
	"sync/atomic"
)

type Service struct {
	Name string
	Replicas []string
}

type Config struct {
	Services []*Service
	Strategy string
}

type Server struct {
	Url *url.URL
	Proxy *httputil.ReverseProxy
}

type ServerList struct {
	Servers []*Server
	Current uint64
}

func (sl *ServerList) NextServer() uint64 {
	next := atomic.AddUint64(&sl.Current, 1)
	serversLen := uint64(len(sl.Servers))
	return next % serversLen
}
