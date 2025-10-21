package logger

import (
	"fmt"
	"log"
	"os"
)

// Logger defines minimal logging capabilities required by the application.
type Logger interface {
	Info(msg string, fields map[string]interface{})
	Error(msg string, err error, fields map[string]interface{})
}

// StdLogger implements Logger using the standard log package.
type StdLogger struct {
	info  *log.Logger
	error *log.Logger
}

// NewStdLogger creates a new StdLogger instance.
func NewStdLogger() *StdLogger {
	return &StdLogger{
		info:  log.New(os.Stdout, "INFO  ", log.LstdFlags|log.Lmsgprefix),
		error: log.New(os.Stderr, "ERROR ", log.LstdFlags|log.Lmsgprefix),
	}
}

// Info logs informational messages.
func (l *StdLogger) Info(msg string, fields map[string]interface{}) {
	l.info.Println(serialize(msg, fields))
}

// Error logs error messages with context.
func (l *StdLogger) Error(msg string, err error, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["error"] = err.Error()
	l.error.Println(serialize(msg, fields))
}

func serialize(msg string, fields map[string]interface{}) string {
	if len(fields) == 0 {
		return msg
	}
	builder := msg
	for key, value := range fields {
		builder += " " + key + "=" + toString(value)
	}
	return builder
}

func toString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int, int64, float64, bool:
		return fmt.Sprintf("%v", v)
	default:
		return fmt.Sprintf("%+v", v)
	}
}
