package logger

import (
	"context"
	"log"

	"go.opentelemetry.io/otel/trace"
)

// Logger ...
type Logger struct {
	*log.Logger
}

type spanLogger struct {
	*log.Logger
	span trace.Span
}

func (sl *spanLogger) Info(msg string) {
	sl.span.AddEvent(msg)
	sl.Logger.Println(msg)
}

func (sl *spanLogger) logToSpan(level string, msg string) {
}

// For ...
func (l *Logger) For(ctx context.Context) *spanLogger {
	if span := trace.SpanFromContext(ctx); span != nil {
		return &spanLogger{
			l.Logger,
			span,
		}
	}

	return nil
}
