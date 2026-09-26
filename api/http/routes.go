package http

import "net/http"

func (s *Server) NewMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /ui", s.getShellHandler)
	mux.HandleFunc("GET /ui/requests", s.getRequestsHandler)
	mux.HandleFunc("GET /ui/requests/{id}", s.getRequestByIDHandler)
	mux.HandleFunc("DELETE /ui/requests", s.resetRequestsHandler)
	mux.HandleFunc("GET /ui/htmx.min.js", s.handleAsset(htmxJS, "text/javascript; charset=utf-8"))
	mux.HandleFunc("GET /ui/styles.css", s.handleAsset(stylesCSS, "text/css; charset=utf-8"))
	mux.HandleFunc("GET /ui/peeko.svg", s.handleAsset(peekoSVG, "image/svg+xml"))
	mux.HandleFunc("GET /ui/peeko-favicon.svg", s.handleAsset(peekoFaviconSVG, "image/svg+xml"))
	mux.HandleFunc("GET /ui/fonts/geist-latin.woff2", s.handleAsset(geistFont, "font/woff2"))

	mux.HandleFunc("/", s.createRequestHandler)

	return mux
}

// Start - starts the HTTP server on the configured address.
func (s *Server) Start() error {
	return http.ListenAndServe(s.addr, s.NewMux())
}
