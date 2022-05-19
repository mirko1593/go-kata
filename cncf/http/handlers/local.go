package handlers

import (
	"bytes"
	"cncf/request"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
)

// HandleLocal ...
func HandleLocal(w http.ResponseWriter, r *http.Request) {
	tr := otel.Tracer("local")
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

}
