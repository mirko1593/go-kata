package handlers

import (
	"cncf/logger"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var l *logger.Logger

func init() {

	l = &logger.Logger{
		Logger: log.New(os.Stdout, "cncf_service:", log.LstdFlags),
	}

}

// HandleFib ...
func HandleFib(w http.ResponseWriter, r *http.Request) {
	tr := otel.Tracer("/fib")
	ctx, span := tr.Start(r.Context(), "Handler")
	defer span.End()
	bs, _ := io.ReadAll(r.Body)
	var data struct {
		Number int64
	}
	json.Unmarshal(bs, &data)

	span.SetAttributes(attribute.Key("http.url").String("/fib"))
	span.SetAttributes(attribute.Key("number").Int64(data.Number))
	span.SetAttributes(attribute.Key("timestamp").String(time.Now().Format("2006-01-02 15:04:05")))

	l.For(ctx).Info("this is a log both to log and span")

	f, _ := func(ctx context.Context) (uint64, error) {
		_, span := otel.Tracer("fib").Start(ctx, "Fibonacci")
		defer span.End()
		return Fibonacci(uint(data.Number))
	}(ctx)

	var response = struct {
		TracerID interface{}
		Data     interface{}
	}{
		TracerID: span.SpanContext().TraceID(),
		Data: map[string]interface{}{
			"Number": f,
		},
	}

	d, _ := json.Marshal(response)

	span.SetAttributes(attribute.Key("http.status_code").Int64(http.StatusOK))
	w.Write(d)
}

// Fibonacci ...
func Fibonacci(n uint) (uint64, error) {
	if n <= 1 {
		return uint64(n), nil
	}

	var n2, n1 uint64 = 0, 1
	for i := uint(2); i < n; i++ {
		n2, n1 = n1, n1+n2
	}

	return n2 + n1, nil
}
