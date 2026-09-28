package webhook

import (
	"crypto/hmac"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Status string

const (
	StatusValid   Status = "valid"
	StatusInvalid Status = "invalid"
	StatusExpired Status = "expired"
	StatusMissing Status = "missing"
	StatusSkipped Status = "skipped"
)

type Encoding string

const (
	EncodingHex    Encoding = "hex"
	EncodingBase64 Encoding = "base64"
)

type Config struct {
	Secret          []byte
	SignatureHeader string        // e.g. "X-Webhook-Signature"
	TimestampHeader string        // "" = body-only signing (no timestamp check)
	Encoding        Encoding      // EncodingHex | EncodingBase64
	Prefix          string        // e.g. "sha256=" or "v1," — includes its separator
	Tolerance       time.Duration // accepted clock skew around the timestamp
}

// Enabled - reports whether verification is active.
func (c Config) Enabled() bool {
	return len(c.Secret) > 0
}

type Verifier struct {
	cfg Config
}

// NewVerifier - creates a verifier from config, applying normalizations.
func NewVerifier(cfg Config) *Verifier {
	if cfg.Encoding == "" {
		cfg.Encoding = EncodingHex
	}

	if cfg.SignatureHeader == "" {
		cfg.SignatureHeader = "X-Webhook-Signature"
	}

	if cfg.Tolerance <= 0 {
		cfg.Tolerance = 5 * time.Minute
	}

	cfg.Prefix = strings.TrimSpace(cfg.Prefix)
	return &Verifier{cfg: cfg}
}

type Result struct {
	Status            Status
	SignedPayload     string // the exact string that was signed
	ReceivedSignature string // as sent on the wire
	ComputedSignature string // computed locally, same encoding
	Note              string // explanation for inconclusive or edge cases
}

// Verify - verifies an HMAC-SHA256 signature over the raw body bytes.
func (v *Verifier) Verify(headers http.Header, body []byte, now time.Time) Result {
	if v == nil || !v.cfg.Enabled() {
		return Result{Status: StatusSkipped, Note: "webhook verification disabled - WEBHOOK_SECRET not set"}
	}

	received := headers.Get(v.cfg.SignatureHeader)
	if received == "" {
		return Result{Status: StatusMissing, Note: fmt.Sprintf("webhook signature header %q not present", v.cfg.SignatureHeader)}
	}

	timestamp := ""
	if v.cfg.TimestampHeader != "" {
		timestamp = headers.Get(v.cfg.TimestampHeader)
		if timestamp == "" {
			return Result{Status: StatusMissing, Note: fmt.Sprintf("webhook timestamp header %q not present", v.cfg.TimestampHeader)}
		}

		ts, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil {
			return Result{Status: StatusExpired, Note: fmt.Sprintf("webhook timestamp %q is not valid unix seconds", timestamp)}
		}

		at := time.Unix(ts, 0)
		if now.Sub(at) > v.cfg.Tolerance || at.Sub(now) > v.cfg.Tolerance {
			return Result{Status: StatusExpired, Note: fmt.Sprintf("webhook timestamp %s outside tolerance %s", timestamp, v.cfg.Tolerance)}
		}
	}

	signature := received
	if v.cfg.Prefix != "" {
		signature = strings.TrimPrefix(signature, v.cfg.Prefix)
	}

	built := v.rebuildSignature(body, timestamp, string(v.cfg.Secret))
	if !hmac.Equal([]byte(built), []byte(signature)) {
		return Result{
			Status:            StatusInvalid,
			SignedPayload:     signedContent(body, timestamp),
			ReceivedSignature: received,
			ComputedSignature: built,
			Note:              "webhook signature mismatch",
		}
	}

	return Result{
		Status:            StatusValid,
		SignedPayload:     signedContent(body, timestamp),
		ReceivedSignature: received,
		ComputedSignature: built,
	}
}
