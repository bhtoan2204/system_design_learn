package logging

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type SpanLogger struct {
	ctx    context.Context
	logger *zap.SugaredLogger
	tracer trace.Tracer
}

func NewSpanLogger(ctx context.Context) *SpanLogger {
	return &SpanLogger{
		ctx:    ctx,
		logger: FromContext(ctx),
		tracer: otel.Tracer("logger-tracer"),
	}
}

func (s *SpanLogger) startSpan(name string) (context.Context, trace.Span) {
	return s.tracer.Start(s.ctx, name)
}

func (s *SpanLogger) Infow(spanName string, msg string, keysAndValues ...interface{}) {
	ctx, span := s.startSpan(spanName)
	defer span.End()

	spanCtx := trace.SpanContextFromContext(ctx)
	fields := []interface{}{
		"trace_id", spanCtx.TraceID().String(),
		"span_id", spanCtx.SpanID().String(),
		"request_id", RequestIDFromContext(ctx),
	}

	s.logger.Infow(msg, append(keysAndValues, fields...)...)
	span.AddEvent(msg, trace.WithAttributes(attribute.String("log.level", "INFO"), attribute.String("log.message", msg)))
}

func (s *SpanLogger) Warnw(spanName string, msg string, err error, keysAndValues ...interface{}) {
	ctx, span := s.startSpan(spanName)
	defer span.End()

	spanCtx := trace.SpanContextFromContext(ctx)
	fields := []interface{}{
		"trace_id", spanCtx.TraceID().String(),
		"span_id", spanCtx.SpanID().String(),
		"request_id", RequestIDFromContext(ctx),
	}

	s.logger.Warnw(msg, append(keysAndValues, fields...)...)
	span.RecordError(err)
	span.AddEvent(msg, trace.WithAttributes(attribute.String("log.level", "WARN"), attribute.String("log.message", msg)))
	span.SetStatus(codes.Error, err.Error())
}

func (s *SpanLogger) Errorw(spanName string, msg string, err error, keysAndValues ...interface{}) {
	ctx, span := s.startSpan(spanName)
	defer span.End()

	spanCtx := trace.SpanContextFromContext(ctx)
	fields := []interface{}{
		"trace_id", spanCtx.TraceID().String(),
		"span_id", spanCtx.SpanID().String(),
		"request_id", RequestIDFromContext(ctx),
	}

	s.logger.Errorw(msg, append(keysAndValues, fields...)...)
	span.RecordError(err)
	span.AddEvent(msg, trace.WithAttributes(attribute.String("log.level", "ERROR"), attribute.String("log.message", msg)))
	span.SetStatus(codes.Error, err.Error())
}
