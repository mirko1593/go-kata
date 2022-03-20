package main

import (
	"context"
	"feature/api"
	"fmt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
)

var (
	addrs = []string{":50051", ":50052"}
)

type ecServer struct {
	api.UnimplementedEchoServer
	addr string
}

func (s *ecServer) UnaryEcho(ctx context.Context, in *api.EchoRequest) (*api.EchoResponse, error) {
	return &api.EchoResponse{Message: fmt.Sprintf("%s (from %s)", in.Message, s.addr)}, nil
}

func startServer(addr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	api.RegisterEchoServer(s, &ecServer{
		addr: addr,
	})
	log.Printf("serving on %v\n", addr)

	s.Serve(lis)
}

func main() {
	var wg sync.WaitGroup
	for _, addr := range addrs {
		wg.Add(1)
		go func(addr string) {
			defer wg.Done()
			startServer(addr)
		}(addr)
	}

	wg.Wait()
}
