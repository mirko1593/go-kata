package handlers

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

var livenessCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "liveness_request_count",
		Help: "No of liveness request",
	},
)

func init() {
	prometheus.MustRegister(livenessCounter)
}

// Liveness ...
func Liveness(w http.ResponseWriter, r *http.Request) {
	livenessCounter.Inc()

	w.Write([]byte("ok"))
}
