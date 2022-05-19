package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/trace"
)

// revive:disable
// HandleOtelServer ...
func HandleOtelServer(w http.ResponseWriter, r *http.Request) {
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
}

// HandleOtelClient ...
func HandleOtelClient(w http.ResponseWriter, r *http.Request) {
	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	// by propagation.Baggage{}
	bag, _ := baggage.Parse("username=donuts")
	ctx := baggage.ContextWithBaggage(r.Context(), bag)

	// by progagation.TracerContext{}
	ctx, span := otel.Tracer("hello service").Start(ctx, "hello client")
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

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost:"+data.Port+"/server/hello", bytes.NewReader(bs))
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

}
