package server

import (
	"context"
	"errors"
	"event_sourcing_bank_system_api/package/logger"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
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
	log := logger.FromContext(ctx)

	errCh := make(chan error, 1)
	go func() {
		<-ctx.Done()
		log.Debug("ServeHTTP: context is closed")

		shutdownCtx, done := context.WithTimeout(context.Background(), 5*time.Second)
		defer done()

		log.Info("ServeHTTP: shutting down")
		errCh <- srv.Shutdown(shutdownCtx)
	}()

	err := srv.Serve(s.listener)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to serve: %w", err)
	}

	log.Info("ServeHTTP: serve stopped")

	err = <-errCh
	return err
}

func (s *Server) ServeHTTPHandler(ctx context.Context, handler http.Handler) error {
	return s.ServeHTTP(ctx, &http.Server{
		ReadHeaderTimeout: 10 * time.Second,
		Handler:           handler,
	})
}

func (s *Server) ServeGRPC(ctx context.Context, srv *grpc.Server) error {
	log := logger.FromContext(ctx)

	go func() {
		<-ctx.Done()

		log.Info("ServeGRPC: shutting down")
		srv.GracefulStop()
	}()

	log.Infow("ServeGRPC: gRPC server started", zap.String("addr", s.listener.Addr().String()))

	if err := srv.Serve(s.listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		return fmt.Errorf("failed to serve: %w", err)
	}

	log.Info("ServeGPRC: serve stopped")

	return nil
}

func (s *Server) Addr() string {
	return net.JoinHostPort(s.ip, s.port)
}

func (s *Server) IP() string {
	return s.ip
}

func (s *Server) Port() string {
	return s.port
}
