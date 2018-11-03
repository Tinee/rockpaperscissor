package http

import (
	"net/http"
	"time"
)

// Server represents the server..
type Server struct {
	*http.Server
}

// NewServer returns a pointer to a Server.
func NewServer(addr string, handler http.Handler) *Server {
	return &Server{
		Server: &http.Server{
			Addr:              addr,
			Handler:           handler,
			MaxHeaderBytes:    1024,
			WriteTimeout:      time.Second * 10,
			ReadTimeout:       time.Second * 10,
			ReadHeaderTimeout: time.Second * 10,
		},
	}
}

// Open opens a connection and start listening on the given port.
func (s *Server) Open() error {
	return s.ListenAndServe()
}

// Close closes the underlaying network listener.
func (s *Server) Close() error {
	return s.Server.Close()
}
