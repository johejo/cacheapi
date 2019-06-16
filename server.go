package cacheapi

import (
	"fmt"
	"net/http"
)

// Server represents cacheapi server.
type Server interface {
	Run() error
}

type server struct {
	host string
	port int
	CacheAPI
}

// NewServer returns a new Server.
func NewServer(api CacheAPI, host string, port int) Server {
	return &server{
		CacheAPI: api,
		host:     host,
		port:     port,
	}
}

// Run listen and serve.
func (s *server) Run() error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	logger.Printf("starting server on: %s", addr)
	return http.ListenAndServe(addr, s.CacheAPI.Handler())
}
