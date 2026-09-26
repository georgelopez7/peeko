package http

import "github.com/georgelopez7/peeko/internal/store"

// Server - the Peeko HTTP server holding the request store.
type Server struct {
	addr  string
	store *store.Store
}

// NewServer - creates a server bound to addr, backed by the given store.
func NewServer(addr string, st *store.Store) *Server {
	return &Server{addr: addr, store: st}
}

// Addr - returns the address the server listens on.
func (s *Server) Addr() string {
	return s.addr
}
