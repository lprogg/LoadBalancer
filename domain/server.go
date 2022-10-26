package domain

import (
	"net/http/httputil"
	"net/url"
)

type Service struct {
	Name string `yml:"name"`
	Matcher string `yml:"matcher"`
	Replicas []string `yml:"replicas"`
	Strategy string `yml:"strategy"`
}

type Server struct {
	Url *url.URL
	Proxy *httputil.ReverseProxy
}
