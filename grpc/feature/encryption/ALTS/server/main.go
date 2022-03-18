package main

import (
	"context"
	"feature/api"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/alts"
)

var port = flag.Int("port", 50051, "the port to serve on")

type ecServer struct {
	api.UnimplementedEchoServer
}

func (s *ecServer) UnaryEcho(ctx context.Context, in *api.EchoRequest) (*api.EchoResponse, error) {
	return &api.EchoResponse{Message: in.Message}, nil
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	altsTC := alts.NewServerCreds(alts.DefaultServerOptions())

	s := grpc.NewServer(grpc.Creds(altsTC))

	api.RegisterEchoServer(s, &ecServer{})

	s.Serve(lis)
}
