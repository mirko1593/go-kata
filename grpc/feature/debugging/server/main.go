package main

import (
	"context"
	"feature/api"
	"log"
	"math/rand"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/channelz/service"
)

var (
	ports = []string{":10001", ":10002", ":10003"}
)

type server struct {
	api.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *api.HelloRequest) (*api.HelloResponse, error) {
	return &api.HelloResponse{Message: "Hello " + in.Name}, nil
}

type slowServer struct {
	api.UnimplementedGreeterServer
}

func (s *slowServer) SayHello(ctx context.Context, in *api.HelloRequest) (*api.HelloResponse, error) {
	time.Sleep(time.Duration(100+rand.Intn(100)) * time.Millisecond)
	return &api.HelloResponse{Message: "Hello " + in.Name}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	s := grpc.NewServer()
	service.RegisterChannelzServiceToServer(s)
	go s.Serve(lis)
	defer s.Stop()

	for i := 0; i < 3; i++ {
		lis, err := net.Listen("tcp", ports[i])
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		defer lis.Close()
		s := grpc.NewServer()
		if i == 2 {
			api.RegisterGreeterServer(s, &slowServer{})
		} else {
			api.RegisterGreeterServer(s, &server{})
		}
		go s.Serve(lis)
	}

	select {}
}
