package domain

import (
	"net/http"
	"sort"
	"strings"
	"time"
)

const MaxBodySize = 64 * 1024

type Request struct {
	ID         int            `json:"id"`
	Method     string         `json:"method"`
	Path       string         `json:"path"`
	RawQuery   string         `json:"raw_query"`
	Host       string         `json:"host"`
	RemoteAddr string         `json:"remote_addr"`
	Proto      string         `json:"proto"`
	Query      map[string]any `json:"query"`
	Headers    map[string]any `json:"headers"`
	Body       string         `json:"body"`
	BodyTrunc  bool           `json:"body_truncated"`
	Length     int            `json:"length"`
	CreatedAt  time.Time      `json:"created_at"`
	Webhook    *WebhookCheck  `json:"webhook,omitempty"`
}

// WebhookCheck - the webhook signature verification outcome for a captured request.
type WebhookCheck struct {
	Status            string `json:"status"`
	SignatureHeader   string `json:"signature_header,omitempty"`
	ReceivedSignature string `json:"received_signature"`
	ComputedSignature string `json:"computed_signature"`
	TimestampHeader   string `json:"timestamp_header,omitempty"`
	ReceivedTimestamp string `json:"received_timestamp,omitempty"`
	Secret            string `json:"secret,omitempty"`
	Encoding          string `json:"encoding,omitempty"`
	Prefix            string `json:"prefix,omitempty"`
	SignedPayload     string `json:"signed_payload"`
	Note              string `json:"note,omitempty"`
}

func NewRequest(id int, r *http.Request, body []byte, truncated bool) Request {
	headers := make(map[string]any, len(r.Header))
	for k, v := range r.Header {
		if k == "Cookie" {
			headers[k] = strings.Join(v, "; ")
			continue
		}
		if len(v) == 1 {
			headers[k] = v[0]
			continue
		}
		headers[k] = v
	}

	query := make(map[string]any)
	for k, v := range r.URL.Query() {
		if len(v) == 1 {
			query[k] = v[0]
			continue
		}
		query[k] = v
	}

	return Request{
		ID:         id,
		Method:     r.Method,
		Path:       r.URL.Path,
		RawQuery:   r.URL.RawQuery,
		Host:       r.Host,
		RemoteAddr: r.RemoteAddr,
		Proto:      r.Proto,
		Query:      query,
		Headers:    headers,
		Body:       string(body),
		BodyTrunc:  truncated,
		Length:     int(r.ContentLength),
		CreatedAt:  time.Now().UTC(),
	}
}

// SortedHeaderKeys - returns the header keys in sorted order for stable rendering.
func (c Request) SortedHeaderKeys() []string {
	keys := make([]string, 0, len(c.Headers))
	for k := range c.Headers {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// SortedQueryKeys - returns the query param keys in sorted order for stable rendering.
func (c Request) SortedQueryKeys() []string {
	keys := make([]string, 0, len(c.Query))
	for k := range c.Query {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
