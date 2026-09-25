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

// htmxJS - the vendored htmx script served from _ui.
var htmxJS, _ = uiFS.ReadFile("_ui/htmx.min.js")

// stylesCSS - the design stylesheet served from _ui.
var stylesCSS, _ = uiFS.ReadFile("_ui/styles.css")

// peekoSVG - the Peeko logo served from _ui.
var peekoSVG, _ = uiFS.ReadFile("_ui/peeko.svg")

// peekoFaviconSVG - the Peeko logo cropped for the favicon, served from _ui.
var peekoFaviconSVG, _ = uiFS.ReadFile("_ui/peeko-favicon.svg")

// geistFont - the Geist variable font (fontsource) served from _ui.
var geistFont, _ = uiFS.ReadFile("_ui/fonts/geist-latin.woff2")

// handleAsset - serves an embedded static asset with the given content type.
func (s *Server) handleAsset(data []byte, contentType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", contentType)
		w.Write(data)
	}
}

// handleShell - serves the htmx index page.
func (s *Server) handleShell(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(shellHTML)
}

// handleList - renders the captured request list fragment; polled by htmx.
func (s *Server) handleList(w http.ResponseWriter, r *http.Request) {
	requests := s.store.List()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := templates.ExecuteTemplate(w, "list", requests); err != nil {
		slog.Error("render list", "err", err)
		http.Error(w, "render error", http.StatusInternalServerError)
	}
}

// handleDetail - renders the detail fragment for a single captured request.
func (s *Server) handleDetail(w http.ResponseWriter, r *http.Request) {
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

// handleClear - drops all captured requests and re-renders the list.
func (s *Server) handleClear(w http.ResponseWriter, r *http.Request) {
	s.store.Clear()
	s.handleList(w, r)
}

// methodClass - maps an HTTP method to its semantic badge class.
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

// shellHTML - the htmx index shell served from _ui.
var shellHTML, _ = uiFS.ReadFile("_ui/index.html")