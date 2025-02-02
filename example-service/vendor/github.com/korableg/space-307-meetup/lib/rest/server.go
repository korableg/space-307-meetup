package rest

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
)

type (
	Handler interface {
		Handler() (string, http.Handler)
	}

	Server struct {
		server *http.Server
	}
)

func New(cfg Config, hdrs ...Handler) *Server {
	var (
		mux    = http.NewServeMux()
		health = new(health)
	)

	mux.Handle("GET /_health", health)
	for _, c := range hdrs {
		mux.Handle(c.Handler())
	}

	srv := &http.Server{
		Addr:        cfg.Address,
		ReadTimeout: cfg.Timeout.Read,
		Handler:     http.TimeoutHandler(mux, cfg.Timeout.Handler, ""),
	}

	srv.RegisterOnShutdown(health.Shutdown)

	return &Server{
		server: srv,
	}
}

func (s *Server) Serve() error {
	slog.Info("Starting HTTP server on " + s.server.Addr)

	lis, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return err
	}

	err = s.server.Serve(lis)
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}

func (s *Server) Shutdown(ctx context.Context) error {
	err := s.server.Shutdown(ctx)
	if err != nil {
		return err
	}

	slog.Info("HTTP server successfully stopped")

	return nil
}
