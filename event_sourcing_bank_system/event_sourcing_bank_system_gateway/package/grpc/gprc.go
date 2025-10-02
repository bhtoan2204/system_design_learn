package grpc

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"

	"go.elastic.co/apm/module/apmgrpc/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	grpcZap "event_sourcing_bank_system_gateway/package/grpc/grpc_zap"
)

const (
	_      = iota //blank identifier
	KB int = 1 << (10 * iota)
	MB
	GB
	TB
	PB
)

// CreateGRPCServer ...
func CreateGRPCServer() *grpc.Server {
	opts := []grpcZap.Option{
		grpcZap.WithDecider(func(fullMethodName string, err error) bool {
			return fullMethodName != "/grpc.health.v1.Health/Check" || err != nil
		}),
	}

	return grpc.NewServer(
		grpc.MaxRecvMsgSize(10*MB),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
			grpcZap.UnaryServerInterceptor(opts...),
		)),
	)
}

// CreateGRPCServerWithRecovery ...
func CreateGRPCServerWithRecovery(f grpc_recovery.RecoveryHandlerFunc) *grpc.Server {
	recoveryOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(f),
	}

	opts := []grpcZap.Option{
		grpcZap.WithDecider(func(fullMethodName string, err error) bool {
			return fullMethodName != "/grpc.health.v1.Health/Check" || err != nil
		}),
	}

	return grpc.NewServer(
		grpc.MaxRecvMsgSize(10*MB),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(recoveryOpts...),
			grpcZap.UnaryServerInterceptor(opts...),
			apmgrpc.NewUnaryServerInterceptor(),
		)),
		grpc.StreamInterceptor(apmgrpc.NewStreamServerInterceptor()),
	)
}

// CreateGRPCClientConn ...
func CreateGRPCClientConn(host string, tlsEnabled bool) (*grpc.ClientConn, error) {
	secureOption := grpc.WithTransportCredentials(insecure.NewCredentials())
	if tlsEnabled {
		creds := credentials.NewTLS(nil)
		secureOption = grpc.WithTransportCredentials(creds)
	}

	return grpc.Dial(
		host,
		secureOption,
		grpc.WithChainUnaryInterceptor(
			grpcZap.UnaryClientInterceptor(),
		),
		// grpc.WithStatsHandler(tracer.GrpcStatsHandler()),
	)
}
