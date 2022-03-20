package main

import (
	"context"
	"feature/api"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

var port = flag.Int("port", 50051, "the port to serve on")

// hwServer is used to implement helloworld.GreeterServer.
type hwServer struct {
	api.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *hwServer) SayHello(ctx context.Context, in *api.HelloRequest) (*api.HelloResponse, error) {
	return &api.HelloResponse{Message: "Hello " + in.Name}, nil
}

type ecServer struct {
	api.UnimplementedEchoServer
}

func (s *ecServer) UnaryEcho(ctx context.Context, req *api.EchoRequest) (*api.EchoResponse, error) {
	return &api.EchoResponse{Message: req.Message}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("server listening at %v\n", lis.Addr())

	s := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	api.RegisterGreeterServer(s, &hwServer{})

	api.RegisterEchoServer(s, &ecServer{})

	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
