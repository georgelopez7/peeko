package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/georgelopez7/peeko/internal/domain"
	"github.com/georgelopez7/peeko/internal/store"
	"github.com/georgelopez7/peeko/internal/webhook"
	"github.com/stretchr/testify/require"
)

const testSecret = "whsec_test"

// testVerifier - a hex verifier signing "<timestamp>.<body>".
func testVerifier() *webhook.Verifier {
	return webhook.NewVerifier(webhook.Config{
		Secret:          []byte(testSecret),
		SignatureHeader: "X-Webhook-Signature",
		TimestampHeader: "X-Webhook-Timestamp",
		Encoding:        webhook.EncodingHex,
		Tolerance:       5 * time.Minute,
	})
}

// signedRequest - builds a POST with a signature over ts + "." + body.
func signedRequest(body string, ts int64) *http.Request {
	req := httptest.NewRequest("POST", "/hook", strings.NewReader(body))
	mac := hmac.New(sha256.New, []byte(testSecret))
	mac.Write([]byte(strconv.FormatInt(ts, 10) + "." + body))
	req.Header.Set("X-Webhook-Signature", hex.EncodeToString(mac.Sum(nil)))
	req.Header.Set("X-Webhook-Timestamp", strconv.FormatInt(ts, 10))
	return req
}

func TestService_InsertRequest(t *testing.T) {
	t.Run("should capture and store a plain request", func(t *testing.T) {
		svc := New(store.NewStore(10), nil)

		captured := svc.InsertRequest(httptest.NewRequest("POST", "/hook", strings.NewReader("hello")))

		require.Equal(t, 0, captured.ID)
		require.Equal(t, "hello", captured.Body)
		require.Nil(t, captured.Webhook, "no verifier - no annotation")

		stored, ok := svc.GetRequestByID(0)
		require.True(t, ok)
		require.Equal(t, "hello", stored.Body)
	})

	t.Run("should annotate a valid signature", func(t *testing.T) {
		svc := New(store.NewStore(10), testVerifier())

		captured := svc.InsertRequest(signedRequest(`{"event":"ping"}`, time.Now().Unix()))

		require.NotNil(t, captured.Webhook)
		require.Equal(t, "valid", captured.Webhook.Status)
		require.Equal(t, captured.Webhook.ReceivedSignature, captured.Webhook.ComputedSignature)
	})

	t.Run("should annotate a tampered body as invalid", func(t *testing.T) {
		svc := New(store.NewStore(10), testVerifier())

		// Sign one body, deliver another.
		req := signedRequest(`{"event":"ping"}`, time.Now().Unix())
		req.Body = io.NopCloser(strings.NewReader(`{"event":"hacked"}`))
		req.ContentLength = int64(len(`{"event":"hacked"}`))

		captured := svc.InsertRequest(req)

		require.NotNil(t, captured.Webhook)
		require.Equal(t, "invalid", captured.Webhook.Status)
		require.NotEqual(t, captured.Webhook.ReceivedSignature, captured.Webhook.ComputedSignature)
	})

	t.Run("should note truncated bodies", func(t *testing.T) {
		svc := New(store.NewStore(10), testVerifier())

		big := strings.Repeat("x", domain.MaxBodySize+10)
		captured := svc.InsertRequest(signedRequest(big, time.Now().Unix()))

		require.True(t, captured.BodyTrunc)
		require.NotNil(t, captured.Webhook)
		require.Contains(t, captured.Webhook.Note, "truncated")
	})
}

func TestService_ListAndReset(t *testing.T) {
	svc := New(store.NewStore(10), nil)

	svc.InsertRequest(httptest.NewRequest("GET", "/a", nil))
	svc.InsertRequest(httptest.NewRequest("GET", "/b", nil))

	list := svc.ListRequests()
	require.Len(t, list, 2)
	require.Equal(t, "/b", list[0].Path, "newest first")

	svc.ResetRequests()
	require.Empty(t, svc.ListRequests())
}
