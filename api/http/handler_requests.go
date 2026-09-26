package http

import (
	"net/http"

	"github.com/georgelopez7/peeko/internal/domain"
)

// createRequestHandler - captures and stores any unmatched request.
func (s *Server) createRequestHandler(w http.ResponseWriter, r *http.Request) {
	captured := domain.NewRequest(0, r)
	s.store.Add(captured)

	w.WriteHeader(http.StatusOK)
}
