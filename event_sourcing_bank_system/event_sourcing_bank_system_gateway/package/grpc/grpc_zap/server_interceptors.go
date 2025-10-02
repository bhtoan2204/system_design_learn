package grpc_zap

import (
	"context"
	"event_sourcing_bank_system_gateway/package/contxt"
	"event_sourcing_bank_system_gateway/package/logger"
	"path"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// SystemField is used in every client-side log statement.
	SystemField = zap.String("system", "grpc")

	// ServerField is used in every server-side log statement.
	ServerField = zap.String("span.kind", "server")
)

type Validator interface {
	Validate() error
}

// UnaryServerInterceptor returns a new unary server interceptors
func UnaryServerInterceptor(opts ...Option) grpc.UnaryServerInterceptor {
	o := evaluateServerOpt(opts)

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		reqID := uuid.NewString()
		ctx = contxt.WithRequestID(ctx, reqID)

		startTime := time.Now()
		newCtx := newLoggerForCall(ctx, info.FullMethod, startTime, reqID)

		// Validate request nếu cần
		if r, ok := req.(Validator); ok {
			if err := r.Validate(); err != nil {
				return nil, status.Error(codes.InvalidArgument, err.Error())
			}
		}

		// Gọi handler
		resp, err := handler(newCtx, req)
		code := o.codeFunc(err)
		if !o.shouldLog(info.FullMethod, err) {
			return resp, err
		}

		duration := time.Since(startTime)
		durField := o.durationFunc(duration)
		level := o.levelFunc(code)
		msg := "finished unary call"

		// Gọi log
		o.messageFunc(newCtx, msg, level, code, err, durField)

		return resp, err
	}
}

func serverCallFields(fullMethodString string) []zapcore.Field {
	service := path.Dir(fullMethodString)[1:]
	method := path.Base(fullMethodString)
	return []zapcore.Field{
		zap.String("system", "grpc"),
		zap.String("span.kind", "server"),
		zap.String("grpc.service", service),
		zap.String("grpc.method", method),
	}
}

func newLoggerForCall(
	ctx context.Context,
	fullMethodString string,
	start time.Time,
	reqID string,
) context.Context {
	log := logger.FromContext(ctx)

	var f []zapcore.Field
	f = append(f, zap.String("grpc.start_time", start.Format(time.RFC3339)))
	f = append(f, zap.String("request_id", reqID))
	if d, ok := ctx.Deadline(); ok {
		f = append(f, zap.String("grpc.request.deadline", d.Format(time.RFC3339)))
	}
	f = append(f, serverCallFields(fullMethodString)...)

	callLog := log.With(f)
	ctx = logger.WithLogger(ctx, callLog)

	return ctx
}
