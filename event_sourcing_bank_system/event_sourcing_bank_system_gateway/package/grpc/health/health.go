package health

import (
	"context"

	pb "google.golang.org/grpc/health/grpc_health_v1"
)

type HealthService struct{}

func NewHealthService() *HealthService { return &HealthService{} }

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
