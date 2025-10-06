package grpc_infra

import (
	"context"
	"event_sourcing_bank_system_api/package/ierror"
	"event_sourcing_bank_system_api/proto/payment"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HandleError() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		response, err := handler(ctx, req)
		if err != nil {
			appErr, ok := err.(*ierror.InternalError)
			if !ok {
				return nil, err
			}
			st := status.New(codes.Code(appErr.GrpcCode), err.Error())
			st, err = st.WithDetails(&payment.ErrorResponse{
				RootError: appErr.RootErr.Error(),
				Message:   appErr.Msg,
				HttpCode:  int64(appErr.HttpCode),
				GrpcCode:  int64(appErr.GrpcCode),
			})
			if err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}
			return nil, st.Err()
		}

		return response, nil
	}
}
