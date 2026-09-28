package service

import (
	"io"
	"net/http"
	"time"

	"github.com/georgelopez7/peeko/internal/domain"
	"github.com/georgelopez7/peeko/internal/webhook"
)

// checkWebhook - runs the webhook verifier over the raw request.
func (s *Service) checkWebhook(headers http.Header, body []byte, truncated bool) *domain.WebhookCheck {
	if s.verifier == nil {
		return nil
	}

	res := s.verifier.Verify(headers, body, time.Now())
	if res.Status == webhook.StatusSkipped {
		return nil
	}

	cfg := s.verifier.Config()
	check := &domain.WebhookCheck{
		Status:            string(res.Status),
		SignatureHeader:   res.SignatureHeader,
		ReceivedSignature: res.ReceivedSignature,
		ComputedSignature: res.ComputedSignature,
		TimestampHeader:   res.TimestampHeader,
		ReceivedTimestamp: res.ReceivedTimestamp,
		Secret:            string(cfg.Secret),
		Encoding:          string(cfg.Encoding),
		Prefix:            cfg.Prefix,
		SignedPayload:     res.SignedPayload,
		Note:              res.Note,
	}

	if truncated {
		check.Note = joinNote(check.Note, "body truncated at capture — signature computed over truncated bytes")
	}

	return check
}

// joinNote - appends a note with a separator.
func joinNote(a, b string) string {
	if a == "" { // TODO: Look into this
		return b
	}

	return a + "; " + b
}

// captureBody - reads the body once, up to MaxBodySize+1 bytes.
func captureBody(r *http.Request) ([]byte, bool) {
	if r.Body == nil {
		return nil, false
	}

	buf := make([]byte, domain.MaxBodySize+1)
	n, _ := io.ReadFull(r.Body, buf)
	if n == 0 {
		return nil, false
	}

	return buf[:min(n, domain.MaxBodySize)], n > domain.MaxBodySize
}
