package service

import (
	"net/http"

	"github.com/georgelopez7/peeko/internal/domain"
	"github.com/georgelopez7/peeko/internal/store"
	"github.com/georgelopez7/peeko/internal/webhook"
)

type Service struct {
	store    *store.Store
	verifier *webhook.Verifier
}

func New(st *store.Store, verifier *webhook.Verifier) *Service {
	return &Service{store: st, verifier: verifier}
}

// InsertRequest - captures an inbound request and stores it.
func (s *Service) InsertRequest(r *http.Request) domain.Request {
	body, truncated := captureBody(r)

	captured := domain.NewRequest(0, r, body, truncated)
	captured.Webhook = s.checkWebhook(r.Header, body, truncated)

	return s.store.Add(captured)
}

// GetRequestByID - returns the captured request by ID.
func (s *Service) GetRequestByID(id int) (domain.Request, bool) {
	return s.store.GetByID(id)
}

// ListRequests - returns all captured requests, newest first.
func (s *Service) ListRequests() []domain.Request {
	return s.store.List()
}

// ResetRequests - removes all captured requests.
func (s *Service) ResetRequests() {
	s.store.Reset()
}
