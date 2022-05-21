package main

import (
	"context"
	"fmt"
	"helloworld/api"
	"log"
	"net"
	"net/http"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

var (
	reg = prometheus.NewRegistry()

	grpcMetrics = grpc_prometheus.NewServerMetrics()

	customizedCounterMetric = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "greeter_server_say_hello_method_handle_count",
		Help: "Total number of RPCs handled on the server.",
	}, []string{"name"})
)

func init() {
	reg.MustRegister(grpcMetrics, customizedCounterMetric)
	customizedCounterMetric.WithLabelValues("Test")
}

type greeterService struct {
	api.UnimplementedGreeterServiceServer
}

func newGreeterService() *greeterService {
	return &greeterService{}
}

// SayHello ...
func (s *greeterService) SayHello(ctx context.Context, in *api.HelloRequest) (*api.HelloResponse, error) {
	customizedCounterMetric.WithLabelValues(in.Name).Inc()
	return &api.HelloResponse{Message: "Hello " + in.Name}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9093))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	httpServer := &http.Server{
		Handler: promhttp.HandlerFor(
			reg,
			promhttp.HandlerOpts{},
		),
		Addr: fmt.Sprintf("0.0.0.0:%d", 9092),
	}

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpcMetrics.StreamServerInterceptor()),
		grpc.UnaryInterceptor(grpcMetrics.UnaryServerInterceptor()),
	)

	api.RegisterGreeterServiceServer(grpcServer, &greeterService{})

	grpcMetrics.InitializeMetrics(grpcServer)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal("Unable to start a http server.")
		}
	}()

	log.Fatal(grpcServer.Serve(lis))
}
