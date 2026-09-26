package http

import (
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

//go:embed all:_ui
var uiFS embed.FS

//go:embed templates.html
var tplFS embed.FS

var templates = template.Must(template.New("").Funcs(template.FuncMap{
	"methodClass": methodClass,
	"formatValue": formatValue,
}).ParseFS(tplFS, "templates.html"))

var htmxJS, _ = uiFS.ReadFile("_ui/htmx.min.js")

var stylesCSS, _ = uiFS.ReadFile("_ui/styles.css")

var peekoSVG, _ = uiFS.ReadFile("_ui/peeko.svg")

var peekoFaviconSVG, _ = uiFS.ReadFile("_ui/peeko-favicon.svg")

var geistFont, _ = uiFS.ReadFile("_ui/fonts/geist-latin.woff2")

var shellHTML, _ = uiFS.ReadFile("_ui/index.html")

// handleAsset - serves an embedded static asset with the given content type.
func (s *Server) handleAsset(data []byte, contentType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", contentType)
		w.Write(data)
	}
}

// getShellHandler - serves the htmx index page.
func (s *Server) getShellHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(shellHTML)
}

// getRequestsHandler - renders the requests.
func (s *Server) getRequestsHandler(w http.ResponseWriter, r *http.Request) {
	requests := s.store.List()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := templates.ExecuteTemplate(w, "list", requests); err != nil {
		slog.Error("render list", "err", err)
		http.Error(w, "render error", http.StatusInternalServerError)
	}
}

// getRequestByIDHandler - renders the details of a request by ID.
func (s *Server) getRequestByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	captured, ok := s.store.Get(id)
	if !ok {
		http.Error(w, "request not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := templates.ExecuteTemplate(w, "detail", captured); err != nil {
		slog.Error("render detail", "err", err)
		http.Error(w, "render error", http.StatusInternalServerError)
	}
}

// resetRequestsHandler - resets all captured requests.
func (s *Server) resetRequestsHandler(w http.ResponseWriter, r *http.Request) {
	s.store.Reset()
	s.getRequestsHandler(w, r)
}

// methodClass - maps an http method to specific css styling.
func methodClass(method string) string {
	switch method {
	case http.MethodGet:
		return "badge badge-get"
	case http.MethodPost:
		return "badge badge-post"
	case http.MethodPut:
		return "badge badge-put"
	case http.MethodPatch:
		return "badge badge-patch"
	case http.MethodDelete:
		return "badge badge-delete"
	case http.MethodHead:
		return "badge badge-head"
	case http.MethodOptions:
		return "badge badge-options"
	default:
		return "badge badge-other"
	}
}

// formatValue - formats a single or multi-value header/query entry.
func formatValue(v any) string {
	switch value := v.(type) {
	case string:
		return value
	case []string:
		return strings.Join(value, ", ")
	default:
		return fmt.Sprint(value)
	}
}
