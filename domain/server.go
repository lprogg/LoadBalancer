package domain

import (
	"fmt"
	"net/http/httputil"
	"net/url"
	"strconv"
	"sync"
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
	Metadata map[string]string
	mutex sync.RWMutex
	alive bool
}

type Replica struct {
	URL string `yml:"url"`
	Metadata map[string]string `yml:"metadata"`
}

func (s *Server) GetMetadataOrDefaultString(key, defaultValue string) string {
	value, ok := s.Metadata[key]

	if !ok {
		return defaultValue
	}

	return value
}

func (s *Server) GetMetadataOrDefaultInt(key string, defaultValue int) int {
	stringValue := s.GetMetadataOrDefaultString(key, fmt.Sprintf("%d", defaultValue))
	intValue, err := strconv.Atoi(stringValue)

	if err != nil {
		return defaultValue
	}

	return intValue
}

func (s *Server) SetServerLiveness(newValue bool) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	oldValue := s.alive
	s.alive = newValue

	return oldValue
}

func (s *Server) IsAlive() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.alive
}
