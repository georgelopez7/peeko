package http

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/georgelopez7/peeko/internal/domain"
	"github.com/georgelopez7/peeko/internal/store"
	"github.com/stretchr/testify/require"
)

// mockRequest - builds a request and adds it to the store for testing.
func mockRequest(s *Server, method, path string, body io.Reader) domain.Request {
	return s.store.Add(domain.NewRequest(0, httptest.NewRequest(method, path, body)))
}

func TestHTTP_CreateRequestHandler(t *testing.T) {
	st := store.NewStore(10)
	s := NewServer("", st)
	mux := s.NewMux()

	t.Run("should capture GET request", func(t *testing.T) {
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, httptest.NewRequest("GET", "/hook?a=1&b=x", nil))

		require.Equal(t, http.StatusOK, rec.Code)

		captured := s.store.List()[0]
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

		captured := s.store.List()[0]
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

		captured := s.store.List()[0]
		require.Equal(t, "PUT", captured.Method)
		require.Equal(t, "/hook", captured.Path)
		require.Equal(t, "updated", captured.Body)
		require.NotNil(t, captured.Query["id"])
	})

	t.Run("should capture PATCH request", func(t *testing.T) {
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, httptest.NewRequest("PATCH", "/hook?id=7", strings.NewReader("patched")))

		require.Equal(t, http.StatusOK, rec.Code)

		captured := s.store.List()[0]
		require.Equal(t, "PATCH", captured.Method)
		require.Equal(t, "/hook", captured.Path)
		require.Equal(t, "patched", captured.Body)
		require.NotNil(t, captured.Query["id"])
	})

	t.Run("should capture DELETE request", func(t *testing.T) {
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, httptest.NewRequest("DELETE", "/hook?id=7", nil))

		require.Equal(t, http.StatusOK, rec.Code)

		captured := s.store.List()[0]
		require.Equal(t, "DELETE", captured.Method)
		require.Equal(t, "/hook", captured.Path)
		require.NotNil(t, captured.Query["id"])
	})

	t.Run("should capture HEAD request", func(t *testing.T) {
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, httptest.NewRequest("HEAD", "/hook", nil))

		require.Equal(t, http.StatusOK, rec.Code)

		captured := s.store.List()[0]
		require.Equal(t, "HEAD", captured.Method)
		require.Equal(t, "/hook", captured.Path)
	})

	t.Run("should capture OPTIONS request", func(t *testing.T) {
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, httptest.NewRequest("OPTIONS", "/hook", nil))

		require.Equal(t, http.StatusOK, rec.Code)

		captured := s.store.List()[0]
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

		require.Empty(t, s.store.List())
	})
}

func TestHTTP_GetRequestByIDHandler(t *testing.T) {
	s := NewServer("", store.NewStore(10))
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

func TestClearHandler(t *testing.T) {
	s := NewServer("", store.NewStore(10))
	mockRequest(s, "GET", "/a", nil)
	mux := s.NewMux()

	t.Run("should DELETE clear return 200", func(t *testing.T) {
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, httptest.NewRequest("DELETE", "/ui/requests", nil))
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("should DELETE clear the store", func(t *testing.T) {
		require.Empty(t, s.store.List())
	})
}
