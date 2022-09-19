package hand

import (
	"context"

	"go.opentelemetry.io/otel"
)

// Hand ...
type Hand struct {
}

// Handle ...
func (*Hand) Handle() {
	tr := otel.Tracer("Hand.Handle")
	_, span := tr.Start(context.Background(), "Hand.Handle")
	span.AddEvent("Hand.Handle")
	defer span.End()
}
