package http

import (
	"net/http"

	"github.com/george-lopez/peeko/internal/domain"
)

// handleCapture - records any unmatched request and replies with a plain 200 OK.
func (s *Server) handleCapture(w http.ResponseWriter, r *http.Request) {
	captured := domain.NewCapturedRequest(0, r)
	s.store.Add(captured)

	w.WriteHeader(http.StatusOK)
}
