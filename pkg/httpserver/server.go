// Package httpserver implements HTTP server.
package httpserver

import (
	"context"
	"net/http"
	"time"
)

const (
	_defaultAddr            = ":80"
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultShutdownTimeout = 3 * time.Second
)

// Server -.
type Server struct {
	App *http.Server

	notify          chan error
	shutdownTimeout time.Duration
}

// New -.
func New(opts ...Option) *Server {
	s := &Server{
		App: &http.Server{
			Addr:         _defaultAddr,
			Handler:      http.DefaultServeMux,
			ReadTimeout:  _defaultReadTimeout,
			WriteTimeout: _defaultWriteTimeout,
		},
		notify:          make(chan error, 1),
		shutdownTimeout: _defaultShutdownTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Start -.
func (s *Server) Start() {
	go func() {
		s.notify <- s.App.ListenAndServe()
		close(s.notify)
	}()
}

// Notify -.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown -.
func (s *Server) Shutdown() error {
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer shutdownCancel()

	err := s.App.Shutdown(shutdownCtx)
	if err == context.DeadlineExceeded {
		return s.App.Close()
	}
	if err != nil {
		return err
	}

	return nil
}
