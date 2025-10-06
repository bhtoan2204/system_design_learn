package grpc_infra

import (
	"context"

	pb "google.golang.org/grpc/health/grpc_health_v1"
)

type HealthService struct{}

func NewHealthService() pb.HealthServer { return &HealthService{} }

func (h *HealthService) Check(ctx context.Context, in *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{
		Status: pb.HealthCheckResponse_SERVING,
	}, nil
}

func (h *HealthService) Watch(in *pb.HealthCheckRequest, srv pb.Health_WatchServer) error {
	return srv.Send(&pb.HealthCheckResponse{
		Status: pb.HealthCheckResponse_SERVING,
	})
}

func (h *HealthService) List(ctx context.Context, in *pb.HealthListRequest) (*pb.HealthListResponse, error) {
	return &pb.HealthListResponse{}, nil
}
