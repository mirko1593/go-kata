package main

import (
	"bufio"
	"context"
	"fmt"
	"helloworld/api"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func main() {
	reg := prometheus.NewRegistry()
	grpcMetrics := grpc_prometheus.NewClientMetrics()
	reg.MustRegister(grpcMetrics)

	conn, err := grpc.Dial(
		fmt.Sprintf("localhost:%v", 9093),
		grpc.WithUnaryInterceptor(grpcMetrics.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(grpcMetrics.StreamClientInterceptor()),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	httpServer := &http.Server{
		Handler: promhttp.HandlerFor(
			reg,
			promhttp.HandlerOpts{},
		),
		Addr: fmt.Sprintf("0.0.0.0:%d", 9094),
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal("Unable to start a http server.")
		}
	}()

	client := api.NewGreeterServiceClient(conn)
	fmt.Println("Start to call the method called SayHello every 3 seconds")
	go func() {
		for {
			_, err := client.SayHello(context.Background(), &api.HelloRequest{Name: "Test"})
			if err != nil {
				log.Printf("Calling the SayHello method unsuccessfully. ErrorInfo: %+v", err)
				log.Printf("You should to stop the process")
				return
			}
			time.Sleep(3 * time.Second)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("You can press n or N to stop the process of client")
	for scanner.Scan() {
		if strings.ToLower(scanner.Text()) == "n" {
			os.Exit(0)
		}
	}
}
