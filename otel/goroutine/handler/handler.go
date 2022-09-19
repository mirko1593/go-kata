package handler

import (
	"context"
	"fmt"
	"routine/hand"
	"sync"

	"go.opentelemetry.io/otel"
)

// Handle ...
func Handle(wg sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		tr := otel.Tracer("Handle")
		_, span := tr.Start(context.Background(), "Handle")
		defer span.End()
		span.AddEvent("Handle handling...")
	}()
}

// SayHello ...
func SayHello(v interface{}) {

	vh, ok := v.(*hand.Hand)
	fmt.Printf("ok: %v\n", ok)

	vh.Handle()
}
