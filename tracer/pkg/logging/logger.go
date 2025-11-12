package logging

import (
	"context"
	"os"
	"strings"
	"sync"
	"time"

	"tracer/constant"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	defaultLoggerOnce sync.Once
	defaultLogger     *zap.SugaredLogger
)

const (
	timestamp  = "timestamp"
	level      = "level"
	logger     = "logger"
	caller     = "caller"
	message    = "message"
	stacktrace = "stacktrace"

	encodingConsole = "console"
	encodingJSON    = "json"
)

var outputStderr = []string{"stderr"}

var productionEncoderConfig = zapcore.EncoderConfig{
	TimeKey:        timestamp,
	LevelKey:       level,
	NameKey:        logger,
	CallerKey:      caller,
	MessageKey:     message,
	StacktraceKey:  stacktrace,
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    levelEncoder(),
	EncodeTime:     timeEncoder(),
	EncodeDuration: zapcore.SecondsDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}

var developmentEncoderConfig = zapcore.EncoderConfig{
	TimeKey:        timestamp,
	LevelKey:       level,
	NameKey:        logger,
	CallerKey:      caller,
	FunctionKey:    zapcore.OmitKey,
	MessageKey:     message,
	StacktraceKey:  stacktrace,
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.CapitalLevelEncoder,
	EncodeTime:     zapcore.ISO8601TimeEncoder,
	EncodeDuration: zapcore.StringDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}

func DefaultLogger() *zap.SugaredLogger {
	defaultLoggerOnce.Do(func() {
		env := os.Getenv("ENVIRONMENT")
		logLevel := os.Getenv("LOG_LEVEL")
		defaultLogger = NewLogger(logLevel, env)
	})
	return defaultLogger
}

func NewLogger(logLevel string, env string) *zap.SugaredLogger {
	level := LogLevel(strings.ToUpper(logLevel))
	var config *zap.Config
	switch env {
	case constant.ENV_PRODUCTION:
		config = &zap.Config{
			Level:            zap.NewAtomicLevelAt(toZapLevel(level)),
			Encoding:         encodingJSON,
			EncoderConfig:    productionEncoderConfig,
			OutputPaths:      outputStderr,
			ErrorOutputPaths: outputStderr,
		}
	case constant.ENV_QC:
		config = &zap.Config{
			Level:            zap.NewAtomicLevelAt(toZapLevel(level)),
			Development:      true,
			Encoding:         encodingJSON,
			EncoderConfig:    developmentEncoderConfig,
			OutputPaths:      outputStderr,
			ErrorOutputPaths: outputStderr,
		}
	default:
		config = &zap.Config{
			Level:            zap.NewAtomicLevelAt(toZapLevel(level)),
			Development:      true,
			Encoding:         encodingConsole,
			EncoderConfig:    developmentEncoderConfig,
			OutputPaths:      outputStderr,
			ErrorOutputPaths: outputStderr,
		}
	}

	logger, err := config.Build()
	if err != nil {
		logger = zap.NewNop()
	}

	return logger.Sugar()
}

func WithLogger(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, constant.LoggerKey, logger)
}

func FromContext(ctx context.Context) *zap.SugaredLogger {
	if logger, ok := ctx.Value(constant.LoggerKey).(*zap.SugaredLogger); ok {
		return logger
	}
	return DefaultLogger()
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	if ctx == nil || requestID == "" {
		return ctx
	}
	return context.WithValue(ctx, constant.RequestIDKey, requestID)
}

func RequestIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if requestID, ok := ctx.Value(constant.RequestIDKey).(string); ok {
		return requestID
	}
	return ""
}

// timeEncoder encodes the time as RFC3339 nano.
func timeEncoder() zapcore.TimeEncoder {
	return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(time.RFC3339Nano))
	}
}
