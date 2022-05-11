package request

import "context"

type tracerIDType struct{}

type spanIDType struct{}

var (
	tracerIDKey = &tracerIDType{}
	spanIDKey   = &spanIDType{}
)

// WithTracerID ...
func WithTracerID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, tracerIDKey, id)
}

// TracerIDFromContext ...
func TracerIDFromContext(ctx context.Context) string {
	if v := ctx.Value(tracerIDKey); v != nil {
		return v.(string)
	}

	return ""
}

// WithSpanID ...
func WithSpanID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, spanIDKey, id)
}

// SpanIDFromContext ...
func SpanIDFromContext(ctx context.Context) string {
	if v := ctx.Value(spanIDKey); v != nil {
		return v.(string)
	}

	return ""
}
