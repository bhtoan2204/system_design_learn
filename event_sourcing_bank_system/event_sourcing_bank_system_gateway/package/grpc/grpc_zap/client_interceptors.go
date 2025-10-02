package grpc_zap

import (
	"context"
	"event_sourcing_bank_system_gateway/package/contxt"
	"event_sourcing_bank_system_gateway/package/logger"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		log := logger.FromContext(ctx)

		ctx = metadata.AppendToOutgoingContext(ctx, "x-request-id", contxt.RequestIDFromCtx(ctx))
		from := time.Now()
		log.Info("Starting grpc request", zap.String("method", method))

		defer func() {
			log.Info("GRPC request finished",
				zap.String("method", method),
				zap.Int64("latency", time.Since(from).Milliseconds()))
		}()

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
