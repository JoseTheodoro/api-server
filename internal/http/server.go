package http

import (
	"context"
	"net/http"
)

type Server struct {
	srv *http.Server
}

func NewServer(listerAddr string, h http.Handler) *Server {
	return &Server{
		srv: &http.Server{
			Addr:    listerAddr,
			Handler: h,
		},
	}
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
