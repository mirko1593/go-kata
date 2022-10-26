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

// SpanLogger ...
type SpanLogger struct {
	*log.Logger
	span trace.Span
}

// Info ...
func (sl *SpanLogger) Info(msg string) {
	sl.span.AddEvent(msg)
	sl.Logger.Println(msg)
}

func (sl *SpanLogger) logToSpan(level string, msg string) {
}

// For ...
func (l *Logger) For(ctx context.Context) *SpanLogger {
	if span := trace.SpanFromContext(ctx); span != nil {
		return &SpanLogger{
			l.Logger,
			span,
		}
	}

	return nil
}
