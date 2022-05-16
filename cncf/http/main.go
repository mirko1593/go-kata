package main

import (
	"bytes"
	"cncf/handlers"
	"cncf/request"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"

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

	otelServerHandler := otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		span := trace.SpanFromContext(ctx)
		bag := baggage.FromContext(ctx)
		span.AddEvent(
			"handling hello request...",
			trace.WithAttributes(
				attribute.Key("username").String(
					bag.Member("username").Value(),
				),
			),
		)

		bs, _ := io.ReadAll(r.Body)
		var data struct {
			Name string
		}
		json.Unmarshal(bs, &data)

		var response = struct {
			TracerID interface{}
			Data     interface{}
		}{
			TracerID: span.SpanContext().TraceID(),
			Data: map[string]interface{}{
				"Message": "Hello " + data.Name,
			},
		}

		d, _ := json.Marshal(response)

		w.Write(d)

	}), "Hello Server")

	http.Handle("/helloserver", otelServerHandler)

	otelClientHandler := otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		client := http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		}

		// by propagation.Baggage{}
		bag, _ := baggage.Parse("username=donuts")
		ctx := baggage.ContextWithBaggage(r.Context(), bag)

		// by progagation.TracerContext{}
		ctx, span := tp.Tracer("hello service").Start(ctx, "hello client")
		defer span.End()

		var data struct {
			Port string
			Name string
		}

		bs, _ := io.ReadAll(r.Body)
		json.Unmarshal(bs, &data)

		bs, _ = json.Marshal(map[string]string{
			"name": data.Name,
		})

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost:"+data.Port+"/helloserver", bytes.NewReader(bs))
		if err != nil {
			log.Fatal(err)
		}

		res, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		w.Write(body)

	}), "Hello Client")

	http.Handle("/helloclient", otelClientHandler)

	http.HandleFunc("/remote", func(w http.ResponseWriter, r *http.Request) {
		tr := otel.Tracer("remote")
		ctx, span := tr.Start(r.Context(), "Remote Handler")
		defer span.End()

		bag, _ := baggage.Parse("username=mirkowang")
		ctx = baggage.ContextWithBaggage(ctx, bag)

		var data struct {
			Port   string
			Number int
		}

		bs, _ := io.ReadAll(r.Body)
		json.Unmarshal(bs, &data)

		bs, _ = json.Marshal(map[string]int{
			"number": data.Number,
		})

		req, err := http.NewRequest(http.MethodPost, "http://localhost:"+data.Port+"/fib", bytes.NewReader(bs))
		if err != nil {
			log.Println("NewRequest", err)
			return
		}

		c := &http.Client{
			Transport: &Transport{
				RoundTripper: http.DefaultTransport,
			},
		}
		remoteRsp, err := c.Do(req.WithContext(
			request.WithSpanID(
				request.WithTracerID(
					ctx,
					span.SpanContext().TraceID().String(),
				),
				span.SpanContext().SpanID().String(),
			),
		))
		if err != nil {
			log.Println("POST", err)
			return
		}
		defer remoteRsp.Body.Close()

		bs, _ = io.ReadAll(remoteRsp.Body)

		w.Write(bs)
	})

	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

// Transport ...
type Transport struct {
	http.RoundTripper
}

// RoundTrip ..
func (t *Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	tid := request.TracerIDFromContext(r.Context())
	if tid != "" {
		log.Println("trace ID:", tid)
		r.Header.Add("tracer-id", tid)
	}

	sid := request.SpanIDFromContext(r.Context())
	if sid != "" {
		log.Println("span ID:", sid)
		r.Header.Add("span-id", sid)
	}

	return t.RoundTripper.RoundTrip(r)
}
