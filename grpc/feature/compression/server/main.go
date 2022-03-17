package main

import (
	"compression/api"
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip" // Install the gzip compressor
)

var port = flag.Int("port", 50051, "the port to serve on")

type server struct {
	api.UnimplementedEchoServer
}

func (s *server) UnaryEcho(ctx context.Context, in *api.EchoRequest) (*api.EchoResponse, error) {
	fmt.Printf("UnaryEcho called with message %q\n", in.GetMessage())
	return &api.EchoResponse{Message: in.Message}, nil
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("server listing at %v\n", lis.Addr())

	s := grpc.NewServer()
	api.RegisterEchoServer(s, &server{})
	s.Serve(lis)
}
