package presentation

import (
	"context"
	"event_sourcing_bank_system_api/infras/grpc_infra"
	"event_sourcing_bank_system_api/package/logger"
	"event_sourcing_bank_system_api/package/server"
	grpclayer "event_sourcing_bank_system_api/presentation/grpc_layer"
	"event_sourcing_bank_system_api/proto/payment"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

type App interface {
	Start(ctx context.Context) error
}

type app struct{}

func NewApp(ctx context.Context) (App, error) {
	return &app{}, nil
}

func (a *app) Start(ctx context.Context) error {
	log := logger.FromContext(ctx)
	log.Info("Starting application")
	panicHandler := func(p any) (err error) {
		return status.Errorf(codes.Internal, "%s", p)
	}
	var sopts []grpc.ServerOption
	sopts = append(sopts,
		grpc.ChainUnaryInterceptor(
			grpc_infra.MonitorRequestDuration(nil),
			grpc_infra.Recovery(panicHandler),
			grpc_infra.Timeout(),
			grpc_infra.HandleError(),
		),
	)
	rpcServer := grpc.NewServer(sopts...)

	healthCheck := grpc_infra.NewHealthService()
	grpc_health_v1.RegisterHealthServer(rpcServer, healthCheck)
	payment.RegisterPaymentServiceServer(rpcServer, grpclayer.NewGrpcPresentation())
	grpcServer, err := server.New(9090)
	if err != nil {
		log.Error("Error creating gRPC server", zap.Error(err))
		return err
	}

	return grpcServer.ServeGRPC(ctx, rpcServer)
}
