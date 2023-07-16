package server

import (
	"context"
	"hezzltask5/internal/handler"

	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
	addr       string
	handler    *handler.Handler
}

func NewServer(addr string, h *handler.Handler) *Server {
	return &Server{
		addr:    addr,
		handler: h,
	}
}

func (s *Server) Run() chan error {
	serverErr := make(chan error)
	s.httpServer = &http.Server{
		Addr:    s.addr,
		Handler: s.handler.Router,
	}

	go func() {
		if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	return serverErr
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}
