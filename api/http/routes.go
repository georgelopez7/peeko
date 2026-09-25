package http

import "net/http"

// NewMux - builds the route table: UI routes plus the catch-all capture handler.
func (s *Server) NewMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /ui", s.handleShell)
	mux.HandleFunc("GET /ui/requests", s.handleList)
	mux.HandleFunc("GET /ui/requests/{id}", s.handleDetail)
	mux.HandleFunc("DELETE /ui/requests", s.handleClear)
	mux.HandleFunc("GET /ui/htmx.min.js", s.handleAsset(htmxJS, "text/javascript; charset=utf-8"))
	mux.HandleFunc("GET /ui/styles.css", s.handleAsset(stylesCSS, "text/css; charset=utf-8"))
	mux.HandleFunc("GET /ui/peeko.svg", s.handleAsset(peekoSVG, "image/svg+xml"))
	mux.HandleFunc("GET /ui/peeko-favicon.svg", s.handleAsset(peekoFaviconSVG, "image/svg+xml"))
	mux.HandleFunc("GET /ui/fonts/geist-latin.woff2", s.handleAsset(geistFont, "font/woff2"))

	mux.HandleFunc("/", s.handleCapture)

	return mux
}

// Start - starts the HTTP server on the configured address.
func (s *Server) Start() error {
	return http.ListenAndServe(s.addr, s.NewMux())
}
