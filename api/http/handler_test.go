package http

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
	"github.com/georgelopez7/peeko/internal/service"
	"github.com/georgelopez7/peeko/internal/store"
	"github.com/georgelopez7/peeko/internal/webhook"
	"github.com/stretchr/testify/require"
)

// mockRequest - builds a request and inserts it via the service for testing.
func mockRequest(s *Server, method, path string, body io.Reader) domain.Request {
	return s.service.InsertRequest(httptest.NewRequest(method, path, body))
}

func TestHTTP_CreateRequestHandler(t *testing.T) {
	st := store.NewStore(10)
	s := NewServer("", service.New(st, nil))
	mux := s.NewMux()

	t.Run("should capture GET request", func(t *testing.T) {
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, httptest.NewRequest("GET", "/hook?a=1&b=x", nil))

		require.Equal(t, http.StatusOK, rec.Code)

		captured := s.service.ListRequests()[0]
		require.Equal(t, "GET", captured.Method)
		require.Equal(t, "/hook", captured.Path)
		require.Empty(t, captured.Body)
		require.NotNil(t, captured.Query["a"])
		require.NotNil(t, captured.Query["b"])
	})

	t.Run("should capture POST request", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/hook?a=1&b=x&b=y", strings.NewReader("hello peeko"))
		req.Header.Set("X-Custom", "yes")
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		captured := s.service.ListRequests()[0]
		require.Equal(t, "POST", captured.Method)
		require.Equal(t, "/hook", captured.Path)
		require.Equal(t, "hello peeko", captured.Body)
		require.Equal(t, "yes", captured.Headers["X-Custom"])
		require.NotNil(t, captured.Query["b"])
	})

	t.Run("should capture PUT request", func(t *testing.T) {
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, httptest.NewRequest("PUT", "/hook?id=7", strings.NewReader("updated")))

		require.Equal(t, http.StatusOK, rec.Code)

		captured := s.service.ListRequests()[0]
		require.Equal(t, "PUT", captured.Method)
		require.Equal(t, "/hook", captured.Path)
		require.Equal(t, "updated", captured.Body)
		require.NotNil(t, captured.Query["id"])
	})

	t.Run("should capture PATCH request", func(t *testing.T) {
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, httptest.NewRequest("PATCH", "/hook?id=7", strings.NewReader("patched")))

		require.Equal(t, http.StatusOK, rec.Code)

		captured := s.service.ListRequests()[0]
		require.Equal(t, "PATCH", captured.Method)
		require.Equal(t, "/hook", captured.Path)
		require.Equal(t, "patched", captured.Body)
		require.NotNil(t, captured.Query["id"])
	})

	t.Run("should capture DELETE request", func(t *testing.T) {
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, httptest.NewRequest("DELETE", "/hook?id=7", nil))

		require.Equal(t, http.StatusOK, rec.Code)

		captured := s.service.ListRequests()[0]
		require.Equal(t, "DELETE", captured.Method)
		require.Equal(t, "/hook", captured.Path)
		require.NotNil(t, captured.Query["id"])
	})

	t.Run("should capture HEAD request", func(t *testing.T) {
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, httptest.NewRequest("HEAD", "/hook", nil))

		require.Equal(t, http.StatusOK, rec.Code)

		captured := s.service.ListRequests()[0]
		require.Equal(t, "HEAD", captured.Method)
		require.Equal(t, "/hook", captured.Path)
	})

	t.Run("should capture OPTIONS request", func(t *testing.T) {
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, httptest.NewRequest("OPTIONS", "/hook", nil))

		require.Equal(t, http.StatusOK, rec.Code)

		captured := s.service.ListRequests()[0]
		require.Equal(t, "OPTIONS", captured.Method)
		require.Equal(t, "/hook", captured.Path)
	})

	t.Run("should serve GET UI routes - no capture", func(t *testing.T) {
		st.Reset()

		for _, path := range []string{"/ui", "/ui/requests", "/ui/styles.css", "/ui/htmx.min.js"} {
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, httptest.NewRequest("GET", path, nil))
			require.Equal(t, http.StatusOK, rec.Code, path)
		}

		require.Empty(t, s.service.ListRequests())
	})
}

func TestHTTP_GetRequestByIDHandler(t *testing.T) {
	s := NewServer("", service.New(store.NewStore(10), nil))
	requestGET := mockRequest(s, "GET", "/thing?a=1", nil)
	requestPOST := mockRequest(s, "POST", "/submit", strings.NewReader("payload"))

	mux := s.NewMux()

	t.Run("should render request by ID - GET", func(t *testing.T) {
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, httptest.NewRequest("GET", "/ui/requests/"+strconv.Itoa(requestGET.ID), nil))
		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Body.String(), "badge-get")
		require.Contains(t, rec.Body.String(), "/thing")
	})

	t.Run("should render request by ID - POST", func(t *testing.T) {
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, httptest.NewRequest("GET", "/ui/requests/"+strconv.Itoa(requestPOST.ID), nil))
		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Body.String(), "badge-post")
		require.Contains(t, rec.Body.String(), "payload")
	})

	t.Run("should return 404 for unknown ID", func(t *testing.T) {
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, httptest.NewRequest("GET", "/ui/requests/999", nil))
		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("should return 400 for non-numeric ID", func(t *testing.T) {
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, httptest.NewRequest("GET", "/ui/requests/notanumber", nil))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestHTTP_ResetRequestsHandler(t *testing.T) {
	s := NewServer("", service.New(store.NewStore(10), nil))
	mockRequest(s, "GET", "/a", nil)
	mux := s.NewMux()

	t.Run("should reset requests", func(t *testing.T) {
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, httptest.NewRequest("DELETE", "/ui/requests", nil))
		require.Equal(t, http.StatusOK, rec.Code)

		require.Empty(t, s.service.ListRequests())
	})
}

func TestHTTP_WebhookVerification(t *testing.T) {
	const secret = "whsec_test"

	newServerWithVerifier := func() *Server {
		cfg := webhook.Config{
			Secret:          []byte(secret),
			SignatureHeader: "X-Webhook-Signature",
			TimestampHeader: "X-Webhook-Timestamp",
			Encoding:        webhook.EncodingHex,
			Tolerance:       5 * time.Minute,
		}
		return NewServer("", service.New(store.NewStore(10), webhook.NewVerifier(cfg)))
	}

	// signedRequest - builds a POST request whose signature covers ts + "." + body.
	signedRequest := func(body string, ts int64) *http.Request {
		req := httptest.NewRequest("POST", "/hook", strings.NewReader(body))
		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write([]byte(strconv.FormatInt(ts, 10) + "." + body))
		req.Header.Set("X-Webhook-Signature", hex.EncodeToString(mac.Sum(nil)))
		req.Header.Set("X-Webhook-Timestamp", strconv.FormatInt(ts, 10))
		return req
	}

	t.Run("should annotate valid signature", func(t *testing.T) {
		s := newServerWithVerifier()
		mux := s.NewMux()

		req := signedRequest(`{"event":"ping"}`, time.Now().Unix())
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		captured := s.service.ListRequests()[0]
		require.NotNil(t, captured.Webhook)
		require.Equal(t, "valid", captured.Webhook.Status)
		require.Equal(t, captured.Webhook.ReceivedSignature, captured.Webhook.ComputedSignature)

		rec = httptest.NewRecorder()
		mux.ServeHTTP(rec, httptest.NewRequest("GET", "/ui/requests/"+strconv.Itoa(captured.ID), nil))
		require.Contains(t, rec.Body.String(), "badge-valid")
	})

	t.Run("should annotate tampered body as invalid", func(t *testing.T) {
		s := newServerWithVerifier()
		mux := s.NewMux()

		req := signedRequest(`{"event":"ping"}`, time.Now().Unix())
		req.Body = io.NopCloser(strings.NewReader(`{"event":"hacked"}`))
		req.ContentLength = int64(len(`{"event":"hacked"}`))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		captured := s.service.ListRequests()[0]
		require.NotNil(t, captured.Webhook)
		require.Equal(t, "invalid", captured.Webhook.Status)
		require.NotEqual(t, captured.Webhook.ReceivedSignature, captured.Webhook.ComputedSignature)

		rec = httptest.NewRecorder()
		mux.ServeHTTP(rec, httptest.NewRequest("GET", "/ui/requests/"+strconv.Itoa(captured.ID), nil))
		require.Contains(t, rec.Body.String(), "badge-invalid")
	})

	t.Run("should annotate missing signature", func(t *testing.T) {
		s := newServerWithVerifier()
		mux := s.NewMux()

		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, httptest.NewRequest("POST", "/hook", strings.NewReader("plain")))

		captured := s.service.ListRequests()[0]
		require.NotNil(t, captured.Webhook)
		require.Equal(t, "missing", captured.Webhook.Status)
	})

	t.Run("should not annotate without verifier", func(t *testing.T) {
		s := NewServer("", service.New(store.NewStore(10), nil))
		mux := s.NewMux()

		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, httptest.NewRequest("POST", "/hook", strings.NewReader("plain")))

		captured := s.service.ListRequests()[0]
		require.Nil(t, captured.Webhook)
	})

	t.Run("should note truncated body as inconclusive", func(t *testing.T) {
		s := newServerWithVerifier()
		mux := s.NewMux()

		// Sign the full oversized body, but Peeko will truncate it at MaxBodySize.
		big := strings.Repeat("x", domain.MaxBodySize+10)
		req := signedRequest(big, time.Now().Unix())
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		captured := s.service.ListRequests()[0]
		require.NotNil(t, captured.Webhook)
		require.True(t, captured.BodyTrunc)
		require.Equal(t, "invalid", captured.Webhook.Status)
		require.Contains(t, captured.Webhook.Note, "truncated")
	})
}
