package domain

import (
	"net/http/httputil"
	"net/url"
)

type Service struct {
	Name string `yml:"name"`
	Matcher string `yml:"matcher"`
	Replicas []Replica `yml:"replicas"`
	Strategy string `yml:"strategy"`
}

type Server struct {
	URL *url.URL
	Proxy *httputil.ReverseProxy
}

type Replica struct {
	URL string `yml:"url"`
	Metadata map[string]string `yml:"metadata"`
}
