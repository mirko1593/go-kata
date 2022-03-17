package main

import (
	"context"
	"feature/api"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
)

var port = flag.Int("port", 50051, "port number")

type server struct {
	api.UnimplementedGreeterServer
	mu    sync.Mutex
	count map[string]int
}

func (s *server) SayHello(ctx context.Context, in *api.HelloRequest) (*api.HelloResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.count[in.Name]++

	if s.count[in.Name] > 1 {
		st := status.New(codes.ResourceExhausted, "Request limit exceeded.")
		ds, err := st.WithDetails(
			&epb.QuotaFailure{
				Violations: []*epb.QuotaFailure_Violation{{
					Subject:     fmt.Sprintf("name:%s", in.Name),
					Description: "Limit one greeting per person",
				}},
			},
		)
		if err != nil {
			return nil, st.Err()
		}

		return nil, ds.Err()
	}

	return &api.HelloResponse{Message: "Hello " + in.Name}, nil
}

func main() {
	flag.Parse()

	address := fmt.Sprintf(":%v", *port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	api.RegisterGreeterServer(s, &server{count: make(map[string]int)})
	s.Serve(lis)
}
