package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"golang.org/x/exp/slog"
)

type Server struct {
	*http.Server
	wg sync.WaitGroup
}

func New(port int, routes http.Handler) *Server {
	return &Server{
		&http.Server{
			Addr:         fmt.Sprintf(":%d", port),
			Handler:      routes,
			IdleTimeout:  1 * time.Minute,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
		sync.WaitGroup{},
	}
}

func (s *Server) Serve() error {
	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sig := <-quit

		slog.Info("caught signal", "signal", sig.String())

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := s.Server.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		slog.Info("completing background tasks", "address", s.Server.Addr)

		s.wg.Wait()
		shutdownError <- nil
	}()

	slog.Info("starting server", "address", s.Server.Addr)

	err := s.Server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	slog.Info("server stopped", "address", s.Server.Addr)

	return nil
}

func (s *Server) Background(fn func()) {
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()

		defer func() {
			if err := recover(); err != nil {
				slog.Error("background recovery error", fmt.Errorf("%s", err)) //nolint:goerr113
			}
		}()

		fn()
	}()
}
