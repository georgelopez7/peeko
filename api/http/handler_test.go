package http

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/george-lopez/peeko/internal/domain"
	"github.com/george-lopez/peeko/internal/store"
)

// capture - builds and stores a captured request from a method/path pair.
func capture(s *Server, method, path string) domain.CapturedRequest {
	return s.store.Add(domain.NewCapturedRequest(0, httptest.NewRequest(method, path, nil)))
}

func TestCaptureHandler(t *testing.T) {
	s := NewServer("", store.NewStore(10))
	req := httptest.NewRequest("POST", "/hook?a=1&b=x&b=y", strings.NewReader("hello peeko"))
	req.Header.Set("X-Custom", "yes")
	rec := httptest.NewRecorder()

	s.NewMux().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rec.Code)
	}

	list := s.store.List()
	if len(list) != 1 {
		t.Fatalf("want 1 captured request, got %d", len(list))
	}

	captured := list[0]
	if captured.Method != "POST" || captured.Path != "/hook" {
		t.Fatalf("bad capture: %s %s", captured.Method, captured.Path)
	}
	if captured.Body != "hello peeko" {
		t.Fatalf("want body captured, got %q", captured.Body)
	}
	if captured.Headers["X-Custom"] != "yes" {
		t.Fatalf("want X-Custom header captured, got %v", captured.Headers["X-Custom"])
	}
	if captured.Query["b"] == nil {
		t.Fatalf("want multi-value query param captured, got %v", captured.Query)
	}
}

func TestUITrafficNotCaptured(t *testing.T) {
	s := NewServer("", store.NewStore(10))
	mux := s.NewMux()

	for _, path := range []string{"/ui", "/ui/requests", "/ui/styles.css", "/ui/htmx.min.js"} {
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, httptest.NewRequest("GET", path, nil))
		if rec.Code != http.StatusOK {
			t.Fatalf("want 200 for %s, got %d", path, rec.Code)
		}
	}

	if got := s.store.List(); len(got) != 0 {
		t.Fatalf("want no captured UI traffic, got %d", len(got))
	}
}

func TestDetailHandler(t *testing.T) {
	s := NewServer("", store.NewStore(10))
	captured := capture(s, "GET", "/thing?a=1")
	mux := s.NewMux()

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("GET", "/ui/requests/"+strconv.Itoa(captured.ID), nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "badge-get") || !strings.Contains(rec.Body.String(), "/thing") {
		t.Fatalf("want detail content, got %q", rec.Body.String())
	}

	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("GET", "/ui/requests/999", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("want 404, got %d", rec.Code)
	}

	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("GET", "/ui/requests/notanumber", nil))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", rec.Code)
	}
}

func TestClearHandler(t *testing.T) {
	s := NewServer("", store.NewStore(10))
	capture(s, "GET", "/a")
	mux := s.NewMux()

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("DELETE", "/ui/requests", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rec.Code)
	}
	if got := s.store.List(); len(got) != 0 {
		t.Fatalf("want store cleared, got %d", len(got))
	}
}
