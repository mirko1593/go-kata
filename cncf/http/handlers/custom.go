package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
)

// revive:disable
// HandleCustomClient ...
func HandleCustomClient(w http.ResponseWriter, r *http.Request) {
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

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost:"+data.Port+"/server/hello/custom", bytes.NewReader(bs))
	if err != nil {
		log.Fatal(err)
	}

	span.AddEvent("manual inject")
	otel.GetTextMapPropagator().Inject(ctx, headerCarrier(req.Header))
	client := http.Client{}
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

// HandleCustomServer ...
func HandleCustomServer(w http.ResponseWriter, r *http.Request) {
	ctx := otel.GetTextMapPropagator().Extract(r.Context(), headerCarrier(r.Header))

	ctx, span := otel.Tracer("Custom").Start(ctx, "HandleCustomServer")
	defer span.End()

	span.AddEvent("this is from custom server handler")

	w.Write([]byte("hello from custom servre"))
}

// Transport ...
type Transport struct {
	http.RoundTripper
}

// RoundTrip ..
func (t *Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	otel.GetTextMapPropagator().Inject(r.Context(), headerCarrier(r.Header))

	// tid := request.TracerIDFromContext(r.Context())
	// if tid != "" {
	// 	log.Println("trace ID:", tid)
	// 	r.Header.Add("tracer-id", tid)
	// }

	// sid := request.SpanIDFromContext(r.Context())
	// if sid != "" {
	// 	log.Println("span ID:", sid)
	// 	r.Header.Add("span-id", sid)
	// }

	return t.RoundTripper.RoundTrip(r)
}

type headerCarrier http.Header

func (hc headerCarrier) Get(key string) string {
	return http.Header(hc).Get(key)
}

func (hc headerCarrier) Set(key, val string) {
	http.Header(hc).Set(key, val)
}

func (hc headerCarrier) Keys() []string {
	keys := make([]string, 0, len(hc))
	for k := range hc {
		keys = append(keys, k)
	}
	return keys
}
