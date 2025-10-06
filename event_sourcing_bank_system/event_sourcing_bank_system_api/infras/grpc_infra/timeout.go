package grpc_infra

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

func Timeout() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ any, err error) {
		ctx, cancel := context.WithTimeout(ctx, time.Second*300)
		defer cancel()

		return handler(ctx, req)
	}
}
