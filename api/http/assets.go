package http

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

//go:embed all:_ui
var uiFS embed.FS

//go:embed templates.html
var tplFS embed.FS

var htmxJS, _ = uiFS.ReadFile("_ui/htmx.min.js")

var stylesCSS, _ = uiFS.ReadFile("_ui/styles.css")

var peekoSVG, _ = uiFS.ReadFile("_ui/assets/peeko.svg")

var geistFont, _ = uiFS.ReadFile("_ui/fonts/geist-latin.woff2")

var geistMonoFont, _ = uiFS.ReadFile("_ui/fonts/geist-mono-latin.woff2")

var shellHTML, _ = uiFS.ReadFile("_ui/index.html")

var templates = template.Must(template.New("").Funcs(template.FuncMap{
	"methodClass": methodClass,
	"statusClass": statusClass,
	"formatValue": formatValue,
	"prettyBody":  prettyBody,
}).ParseFS(tplFS, "templates.html"))

// handleAsset - serves an embedded static asset with the given content type.
func (s *Server) handleAsset(data []byte, contentType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", contentType)
		w.Write(data)
	}
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

// statusClass - maps a webhook status to specific css styling.
func statusClass(status string) string {
	switch status {
	case "valid":
		return "badge badge-valid"
	case "invalid":
		return "badge badge-invalid"
	case "expired":
		return "badge badge-expired"
	case "missing":
		return "badge badge-missing"
	default:
		return "badge badge-skipped"
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

// prettyBody - pretty-prints the request body when it's a form or JSON payload.
func prettyBody(body string) string {
	if strings.HasPrefix(body, "{") || strings.HasPrefix(body, "[") {
		var buf bytes.Buffer
		if err := json.Indent(&buf, []byte(body), "", "  "); err == nil && buf.Len() > 0 {
			return buf.String()
		}
	}

	// form-encoded: "a=1&b=2" -> "a=1\nb=2"
	if strings.Contains(body, "=") && !strings.Contains(body, "\n") {
		lines := strings.Split(body, "&")
		valid := true
		for _, line := range lines {
			if _, _, hasValue := strings.Cut(line, "="); !hasValue {
				valid = false
				break
			}
		}
		if valid {
			return strings.Join(lines, "\n")
		}
	}

	return body
}
