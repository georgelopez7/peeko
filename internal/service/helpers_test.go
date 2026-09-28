package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/georgelopez7/peeko/internal/domain"
	"github.com/georgelopez7/peeko/internal/webhook"
	"github.com/stretchr/testify/require"
)

func Test_captureBody(t *testing.T) {
	t.Run("should return empty for a nil body", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/hook", nil)
		req.Body = nil

		body, truncated := captureBody(req)

		require.Nil(t, body)
		require.False(t, truncated)
	})

	t.Run("should capture a small body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/hook", strings.NewReader("hello"))

		body, truncated := captureBody(req)

		require.Equal(t, "hello", string(body))
		require.False(t, truncated)
	})

	t.Run("should return empty for an empty body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/hook", strings.NewReader(""))

		body, truncated := captureBody(req)

		require.Nil(t, body)
		require.False(t, truncated)
	})

	t.Run("should keep exactly MaxBodySize without truncating", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/hook", strings.NewReader(strings.Repeat("x", domain.MaxBodySize)))

		body, truncated := captureBody(req)

		require.Len(t, body, domain.MaxBodySize)
		require.False(t, truncated)
	})

	t.Run("should truncate an oversized body", func(t *testing.T) {
		big := strings.Repeat("x", domain.MaxBodySize+10)
		req := httptest.NewRequest("POST", "/hook", strings.NewReader(big))

		body, truncated := captureBody(req)

		require.Len(t, body, domain.MaxBodySize)
		require.Equal(t, big[:domain.MaxBodySize], string(body))
		require.True(t, truncated)
	})
}

func Test_checkWebhook(t *testing.T) {
	t.Run("should return nil without a verifier", func(t *testing.T) {
		svc := New(nil, nil)

		check := svc.checkWebhook(nil, []byte("body"), false)

		require.Nil(t, check)
	})

	t.Run("should skip verification for a truncated body", func(t *testing.T) {
		svc := New(nil, testVerifier())

		check := svc.checkWebhook(nil, []byte(strings.Repeat("x", domain.MaxBodySize)), true)

		require.NotNil(t, check)
		require.Equal(t, string(webhook.StatusSkipped), check.Status)
		require.Equal(t, webhook.NoteBodyTruncated, check.Note)
	})

	t.Run("should verify a non-truncated body", func(t *testing.T) {
		svc := New(nil, testVerifier())
		ts := time.Now().Unix()

		body := `{"event":"ping"}`
		mac := hmac.New(sha256.New, []byte(testSecret))
		mac.Write([]byte(strconv.FormatInt(ts, 10) + "." + body))

		h := http.Header{}
		h.Set("X-Webhook-Signature", hex.EncodeToString(mac.Sum(nil)))
		h.Set("X-Webhook-Timestamp", strconv.FormatInt(ts, 10))

		check := svc.checkWebhook(h, []byte(body), false)

		require.NotNil(t, check)
		require.Equal(t, "valid", check.Status)
		require.Equal(t, check.ReceivedSignature, check.ComputedSignature)
	})
}
