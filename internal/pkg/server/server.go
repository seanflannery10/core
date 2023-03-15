package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/exp/slog"
)

func Serve(port int, routes http.Handler) error {
	shutdownError := make(chan error)

	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      routes,
		IdleTimeout:  1 * time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sig := <-quit

		slog.Info("caught signal", "signal", sig.String())

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := s.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		shutdownError <- nil
	}()

	slog.Info("starting server", "address", s.Addr)

	err := s.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed listen and serve: %w", err)
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	slog.Info("server stopped", "address", s.Addr)

	return nil
}
