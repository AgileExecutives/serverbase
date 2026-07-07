package http

import (
	stdhttp "net/http"
)

// Server is a very small HTTP server wrapper that exposes a tiny router surface
// used by modules to register routes.
type Server struct {
	Addr string
	mux  *stdhttp.ServeMux
}

// New creates a new Server listening on addr (e.g. ":8080").
func New(addr string) *Server {
	return &Server{Addr: addr, mux: stdhttp.NewServeMux()}
}

// RegisterRoute registers a handler for the given pattern.
func (s *Server) RegisterRoute(pattern string, h stdhttp.Handler) {
	s.mux.Handle(pattern, h)
}

// Handler returns the underlying http.Handler for testing or composition.
func (s *Server) Handler() stdhttp.Handler { return s.mux }

// Start blocks and serves using the configured address.
func (s *Server) Start() error {
	return stdhttp.ListenAndServe(s.Addr, s.mux)
}
