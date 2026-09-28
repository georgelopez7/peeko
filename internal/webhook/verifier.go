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
	SignatureHeader   string // name of the header the signature was read from
	TimestampHeader   string // name of the header the timestamp was read from, if configured
	ReceivedTimestamp string // raw timestamp header value, if used
	Note              string // explanation for inconclusive or edge cases
}

// Config - returns the configuration the verifier runs with.
func (v *Verifier) Config() Config {
	return v.cfg
}

// Verify - verifies an HMAC-SHA256 signature over the raw body bytes.
func (v *Verifier) Verify(headers http.Header, body []byte, now time.Time) Result {
	if v == nil || !v.cfg.Enabled() {
		return Result{Status: StatusSkipped, Note: "webhook verification disabled - WEBHOOK_SECRET not set"}
	}

	res := Result{SignatureHeader: v.cfg.SignatureHeader}

	received := headers.Get(v.cfg.SignatureHeader)
	if received == "" {
		res.Status = StatusMissing
		res.Note = fmt.Sprintf("webhook signature header %q not present", v.cfg.SignatureHeader)
		return res
	}

	res.ReceivedSignature = received
	res.TimestampHeader = v.cfg.TimestampHeader

	timestamp := ""
	if v.cfg.TimestampHeader != "" {
		timestamp = headers.Get(v.cfg.TimestampHeader)
		res.ReceivedTimestamp = timestamp
		if timestamp == "" {
			res.Status = StatusMissing
			res.Note = fmt.Sprintf("webhook timestamp header %q not present", v.cfg.TimestampHeader)
			return res
		}

		ts, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil {
			res.Status = StatusExpired
			res.Note = fmt.Sprintf("webhook timestamp %q is not valid unix seconds", timestamp)
			return res
		}

		at := time.Unix(ts, 0)
		if now.Sub(at) > v.cfg.Tolerance || at.Sub(now) > v.cfg.Tolerance {
			res.Status = StatusExpired
			res.Note = fmt.Sprintf("webhook timestamp %s outside tolerance %s", timestamp, v.cfg.Tolerance)
			return res
		}
	}

	signature := received
	if v.cfg.Prefix != "" {
		signature = strings.TrimPrefix(signature, v.cfg.Prefix)
	}

	built := v.rebuildSignature(body, timestamp, string(v.cfg.Secret))
	res.ComputedSignature = built
	res.SignedPayload = signedContent(body, timestamp)

	if !hmac.Equal([]byte(built), []byte(signature)) {
		res.Status = StatusInvalid
		res.Note = "webhook signature mismatch"
		return res
	}

	res.Status = StatusValid
	return res
}
