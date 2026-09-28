package http

import (
	"github.com/georgelopez7/peeko/internal/service"
)

type Server struct {
	addr    string
	service *service.Service
}

// NewServer - creates a server bound to addr, backed by the given service.
func NewServer(addr string, svc *service.Service) *Server {
	return &Server{addr: addr, service: svc}
}

// Addr - returns the address the server listens on.
func (s *Server) Addr() string {
	return s.addr
}
