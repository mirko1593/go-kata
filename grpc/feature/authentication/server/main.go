package main

import (
	"context"
	"crypto/tls"
	"feature/api"
	"feature/data"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
)

var port = flag.Int("port", 50051, "the port to serve on")

func main() {
	flag.Parse()
	fmt.Printf("server starting on port %d...\n", *port)

	cert, err := tls.LoadX509KeyPair(data.Path("x509/server_cert.pem"), data.Path("x509/server_key.pem"))
	if err != nil {
		log.Fatalf("failed to load key pair: %s", err)
	}
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(ensureValidToken),
		grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
	}
	s := grpc.NewServer(opts...)
	api.RegisterEchoServer(s, &ecServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type ecServer struct {
	api.UnimplementedEchoServer
}

func (s *ecServer) UnaryEcho(ctx context.Context, in *api.EchoRequest) (*api.EchoResponse, error) {
	return &api.EchoResponse{Message: in.Message}, nil
}

func ensureValidToken(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetadata
	}

	if !valid(md["authorization"]) {
		return nil, errInvalidToken
	}

	return handler(ctx, req)
}

func valid(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}

	token := strings.TrimPrefix(authorization[0], "Bearer ")

	return token == "some-secret-token"
}
