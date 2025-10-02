package grpc

import (
	"context"
	"database/sql"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
)

var (
	ErrNotFound         = errors.New("Not found")
	ErrNoCtxMetaData    = errors.New("No ctx metadata")
	ErrInvalidSessionId = errors.New("Invalid session id")
	ErrEmailExists      = errors.New("Email already exists")
)

// Parse error and get code
func ParseGRPCErrStatusCode(err error) codes.Code {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return codes.NotFound
	// case errors.Is(err, redis.Nil):
	// 	return codes.NotFound
	case errors.Is(err, context.Canceled):
		return codes.Canceled
	case errors.Is(err, context.DeadlineExceeded):
		return codes.DeadlineExceeded
	case errors.Is(err, ErrEmailExists):
		return codes.AlreadyExists
	case errors.Is(err, ErrNoCtxMetaData):
		return codes.Unauthenticated
	case errors.Is(err, ErrInvalidSessionId):
		return codes.PermissionDenied
	case strings.Contains(err.Error(), "Validate"):
		return codes.InvalidArgument
	case strings.Contains(err.Error(), "redis"):
		return codes.NotFound
	}
	return codes.Internal
}

// Map GRPC errors codes to http status
func MapGRPCErrCodeToHttpStatus(code codes.Code) int {
	switch code {
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.AlreadyExists:
		return http.StatusBadRequest
	case codes.NotFound:
		return http.StatusNotFound
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Canceled:
		return http.StatusRequestTimeout
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.InvalidArgument:
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}

// func CatchError() grpc.UnaryServerInterceptor {
// 	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
// 		resp, err := handler(ctx, req)
// 		if err != nil {
// 			code := status.Code(err)
// 			if ctx.Err() != context.Canceled &&
// 				(code == codes.Internal ||
// 					code == codes.Unknown) {
// 				reqID, _ := grpc_zap.FromContext(ctx)
// 				reqJson, _ := json.Marshal(req)
// 				discorde.WithScope(func(s *discorde.Scope) {
// 					s.SetTag("method", info.FullMethod)
// 					s.SetTag("request-id", reqID)
// 					s.SetTag("request-info", string(reqJson))
// 					discorde.CaptureExeption(err)
// 				})
// 			}
// 		}

// 		return resp, err
// 	}
// }
