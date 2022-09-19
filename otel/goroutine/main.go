package main

import (
	"context"
	"fmt"
	"log"
	"routine/hand"
	"routine/handler"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

func tracerProvider(url string) (*tracesdk.TracerProvider, error) {
	exportor, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSyncer(exportor),
		tracesdk.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("test-goroutine"),
			),
		),
	)

	otel.SetTracerProvider(tp)

	return tp, nil
}

func main() {
	tp, err := tracerProvider("http://localhost:14268/api/traces")
	if err != nil {
		log.Fatal(err)
	}

	tr := tp.Tracer("main")
	_, span := tr.Start(context.Background(), "main")
	defer span.End()
	span.AddEvent("Hello World")

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		tr := otel.Tracer("goroutine")
		_, span := tr.Start(context.Background(), "goroutine")
		defer span.End()
		span.AddEvent("goroutine world")
	}()

	handler.Handle(wg)
	handler.SayHello(&hand.Hand{})

	wg.Wait()

	fmt.Println("Done main")
}
