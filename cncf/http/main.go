package main

import (
	"cncf/handlers"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"

	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

var build = "develop"
var pid = os.Getpid()

var port = flag.String("port", "8080", "port for service to listen on")

// ...
const (
	service     = "cncf-service"
	environment = "develop"
)

func tracerProvider(url string) (*tracesdk.TracerProvider, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			attribute.String("environment", environment),
			attribute.Int("PID", pid),
		)),
	)

	return tp, nil
}

func main() {
	flag.Parse()

	tp, err := tracerProvider("http://localhost:14268/api/traces")
	if err != nil {
		log.Fatal(err)
	}

	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	http.HandleFunc("/liveness", handlers.Liveness)

	http.HandleFunc("/fib", handlers.HandleFib)

	http.Handle("/server/hello", otelhttp.NewHandler(
		http.HandlerFunc(handlers.HandleHelloServer),
		"Hello-Server",
	))

	http.Handle("/client/hello/otel", otelhttp.NewHandler(
		http.HandlerFunc(handlers.HandleHelloClient),
		"Hello-Client-Otel",
	))

	// manual inject and extract
	http.HandleFunc("/client/hello/custom", handlers.HandleCustomClient)

	http.HandleFunc("/server/hello/custom", handlers.HandleCustomServer)

	http.HandleFunc("/local", handlers.HandleLocal)

	// scrape route for prometheus
	http.Handle("/metrics", promhttp.Handler())

	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
