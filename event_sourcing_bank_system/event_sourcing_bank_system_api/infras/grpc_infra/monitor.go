package grpc_infra

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/grpc"
)

func MonitorRequestDuration(trace metric.Float64Histogram) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		if err != nil {
			return resp, err
		}

		duration := time.Since(start)

		if trace != nil && shouldMonitor(info) {
			trace.Record(ctx, float64(duration.Milliseconds()), metric.WithAttributes(
				attribute.Key("path").String(info.FullMethod),
			))
		}

		return resp, err
	}
}

func shouldMonitor(info *grpc.UnaryServerInfo) bool {
	return info.FullMethod != "/grpc.health.v1.Health/Check"
}
