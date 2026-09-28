package http

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/georgelopez7/peeko/internal/service"
)

type Server struct {
	addr    string
	srv     *http.Server
	service *service.Service
}

// NewServer - creates a server bound to addr, backed by the given service.
func NewServer(addr string, svc *service.Service) *Server {
	return &Server{addr: addr, srv: &http.Server{Addr: addr}, service: svc}
}

// Addr - returns the address the server listens on.
func (s *Server) Addr() string {
	return s.addr
}

// Start - starts the HTTP server on the configured address.
func (s *Server) Start() error {
	s.srv.Handler = s.NewMux()
	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

// Stop - gracefully shuts down the server, waiting for in-flight requests.
func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.srv.Shutdown(ctx)
}
