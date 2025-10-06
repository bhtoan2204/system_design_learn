package grpclayer

import (
	"event_sourcing_bank_system_api/proto/payment"

	"google.golang.org/grpc"
)

var _ GrpcPresentation = (*grpcPresentation)(nil)

type GrpcPresentation interface {
	payment.PaymentServiceServer
}

type grpcPresentation struct {
	server *grpc.Server
}

func NewGrpcPresentation() GrpcPresentation {
	return &grpcPresentation{
		server: grpc.NewServer(),
	}
}

func (p *grpcPresentation) GRPCServer() *grpc.Server {
	return p.server
}
