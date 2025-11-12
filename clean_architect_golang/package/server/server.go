package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	ip       string
	port     string
	listener net.Listener
}

func New(port int) (*Server, error) {
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &Server{
		ip:       listener.Addr().(*net.TCPAddr).IP.String(),
		port:     strconv.Itoa(listener.Addr().(*net.TCPAddr).Port),
		listener: listener,
	}, nil
}

func (s *Server) ServeHTTP(ctx context.Context, srv *http.Server) error {
	errCh := make(chan error, 1)
	go func() {
		<-ctx.Done()
		log.Println("DEBUG: ServeHTTP: context is closed")

		shutdownCtx, done := context.WithTimeout(context.Background(), 5*time.Second)
		defer done()

		log.Println("INFO: ServeHTTP: shutting down")
		errCh <- srv.Shutdown(shutdownCtx)
	}()

	// This will block until the context is closed.
	err := srv.Serve(s.listener)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to serve: %w", err)
	}

	log.Println("INFO: ServeHTTP: serve stopped")

	err = <-errCh
	return err
}

func (s *Server) ServeHTTPHandler(ctx context.Context, handler http.Handler) error {
	return s.ServeHTTP(ctx, &http.Server{
		ReadHeaderTimeout: 30 * time.Second,
		Handler:           handler,
	})
}
