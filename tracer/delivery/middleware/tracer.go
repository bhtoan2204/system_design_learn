package middleware

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"tracer/pkg/logging"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var httpTracer = otel.Tracer("tracer/delivery/http")

// TraceRequest creates a span for every HTTP request and seeds the trace ID with the request ID.
func TraceRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		requestID := logging.RequestIDFromContext(ctx)

		route := c.FullPath()
		if route == "" {
			route = c.Request.URL.Path
		}

		spanName := fmt.Sprintf("%s %s", c.Request.Method, route)

		attrs := []attribute.KeyValue{
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.route", route),
			attribute.String("http.target", c.Request.URL.Path),
			attribute.String("http.scheme", httpScheme(c)),
		}
		if requestID != "" {
			attrs = append(attrs, attribute.String("http.request_id", requestID))
		}
		if host := c.Request.Host; host != "" {
			attrs = append(attrs, attribute.String("http.host", host))
		}

		spanOptions := []trace.SpanStartOption{
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(attrs...),
		}

		if spanCtx, ok := spanContextFromRequestID(requestID); ok {
			spanOptions = append(spanOptions, trace.WithLinks(trace.Link{
				SpanContext: spanCtx,
			}))
		}

		ctx, span := httpTracer.Start(ctx, spanName, spanOptions...)
		defer span.End()

		c.Request = c.Request.WithContext(ctx)

		c.Next()

		statusCode := c.Writer.Status()
		span.SetAttributes(attribute.Int("http.status_code", statusCode))
		if statusCode >= http.StatusBadRequest {
			span.SetStatus(codes.Error, http.StatusText(statusCode))
		} else {
			span.SetStatus(codes.Ok, http.StatusText(statusCode))
		}
	}
}

func spanContextFromRequestID(requestID string) (trace.SpanContext, bool) {
	traceID, ok := traceIDFromRequestID(requestID)
	if !ok {
		return trace.SpanContext{}, false
	}

	spanID := spanIDFromRequestID(requestID)
	if !spanID.IsValid() {
		return trace.SpanContext{}, false
	}

	spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: trace.FlagsSampled,
		Remote:     true,
	})

	return spanCtx, spanCtx.IsValid()
}

func traceIDFromRequestID(requestID string) (trace.TraceID, bool) {
	if requestID == "" {
		return trace.TraceID{}, false
	}

	uid, err := uuid.Parse(requestID)
	if err != nil {
		return trace.TraceID{}, false
	}

	var traceID trace.TraceID
	copy(traceID[:], uid[:])

	return traceID, traceID.IsValid()
}

func spanIDFromRequestID(requestID string) trace.SpanID {
	sum := sha256.Sum256([]byte(requestID))
	var spanID trace.SpanID
	copy(spanID[:], sum[:len(spanID)])
	return spanID
}

func httpScheme(c *gin.Context) string {
	if c.Request.TLS != nil {
		return "https"
	}
	if c.Request.URL != nil && c.Request.URL.Scheme != "" {
		return c.Request.URL.Scheme
	}
	return "http"
}
