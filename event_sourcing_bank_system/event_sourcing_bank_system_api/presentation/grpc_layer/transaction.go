package grpclayer

import (
	"context"
	"event_sourcing_bank_system_api/package/logger"
	"event_sourcing_bank_system_api/proto/payment"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (p *grpcPresentation) CreateTransaction(ctx context.Context, req *payment.CreateTransactionRequest) (*payment.CreateTransactionResponse, error) {
	log := logger.FromContext(ctx)
	log.Infow("CreateTransaction", zap.Any("req", req))
	return nil, status.Errorf(codes.Unimplemented, "method CreateTransaction not implemented")
}
