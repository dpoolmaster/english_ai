package httpserver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"
)

// Server wraps http.Server with graceful shutdown driven by context
// cancellation.
type Server struct {
	srv *http.Server
	log *slog.Logger
}

func New(addr string, handler http.Handler, log *slog.Logger) *Server {
	return &Server{
		srv: &http.Server{
			Addr:              addr,
			Handler:           handler,
			ReadHeaderTimeout: 5 * time.Second,
			IdleTimeout:       120 * time.Second,
		},
		log: log,
	}
}

// Run serves until the listener fails or ctx is cancelled. On cancellation it
// stops accepting new connections and waits up to shutdownTimeout for active
// requests to finish.
func (s *Server) Run(ctx context.Context, shutdownTimeout time.Duration) error {
	listener, err := net.Listen("tcp", s.srv.Addr)
	if err != nil {
		return fmt.Errorf("listen %s: %w", s.srv.Addr, err)
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- s.srv.Serve(listener)
	}()

	s.log.Info("server started", "addr", s.srv.Addr)

	select {
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("serve %s: %w", s.srv.Addr, err)
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	s.log.Info("server stopping", "addr", s.srv.Addr)
	if err := s.srv.Shutdown(shutdownCtx); err != nil {
		_ = s.srv.Close()
		return fmt.Errorf("shutdown %s: %w", s.srv.Addr, err)
	}

	if err := <-errCh; err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serve %s: %w", s.srv.Addr, err)
	}
	return nil
}
