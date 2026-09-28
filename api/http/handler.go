package http

import (
	"log/slog"
	"net/http"
	"strconv"
)

// createRequestHandler - captures and stores any unmatched request.
func (s *Server) createRequestHandler(w http.ResponseWriter, r *http.Request) {
	s.service.InsertRequest(r)

	w.WriteHeader(http.StatusOK)
}

// getRequestByIDHandler - renders the details of a request by ID.
func (s *Server) getRequestByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	captured, ok := s.service.GetRequestByID(id)
	if !ok {
		http.Error(w, "request not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := templates.ExecuteTemplate(w, "detail", captured); err != nil {
		slog.Error("render detail", "err", err)
		http.Error(w, "render error", http.StatusInternalServerError)
	}
}

// getRequestsHandler - renders the requests.
func (s *Server) getRequestsHandler(w http.ResponseWriter, r *http.Request) {
	requests := s.service.ListRequests()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := templates.ExecuteTemplate(w, "list", requests); err != nil {
		slog.Error("render list", "err", err)
		http.Error(w, "render error", http.StatusInternalServerError)
	}
}

// resetRequestsHandler - resets all captured requests.
func (s *Server) resetRequestsHandler(w http.ResponseWriter, r *http.Request) {
	s.service.ResetRequests()
	s.getRequestsHandler(w, r)
}
