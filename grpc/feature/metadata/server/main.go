package main

import (
	"context"
	"feature/api"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var port = flag.Int("port", 50051, "the port to serve on")

const (
	timestampFormat = time.StampNano
	streamingCount  = 10
)

type server struct {
	api.UnimplementedEchoServer
}

func (s *server) UnaryEcho(ctx context.Context, in *api.EchoRequest) (*api.EchoResponse, error) {
	fmt.Printf("--- UnaryEcho ---\n")
	defer func() {
		trailer := metadata.Pairs("timestamp", time.Now().Format(timestampFormat))
		grpc.SetTrailer(ctx, trailer)
	}()

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.DataLoss, "UnaryEcho: failed to get metadata")
	}
	if t, ok := md["timestamp"]; ok {
		fmt.Printf("timestamp from metadata: \n")
		for i, e := range t {
			fmt.Printf(" %d. %s\n", i, e)
		}
	}

	// create and send header
	header := metadata.New(map[string]string{"location": "MTV", "timestamp": time.Now().Format(timestampFormat)})
	grpc.SendHeader(ctx, header)

	fmt.Printf("request received: %v, sending echo\n", in)

	return &api.EchoResponse{Message: in.Message}, nil
}

func (s *server) ServerStreamingEcho(in *api.EchoRequest, stream api.Echo_ServerStreamingEchoServer) error {
	fmt.Printf("--- ServerStreamingEcho --- \n")
	defer func() {
		trailer := metadata.Pairs("timestamp", time.Now().Format(timestampFormat))
		stream.SetTrailer(trailer)
	}()

	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return status.Errorf(codes.DataLoss, "ServerStreamingEcho: failed to get metadata")
	}

	if t, ok := md["timestamp"]; ok {
		fmt.Printf("timestamp from metadata: \n")
		for i, e := range t {
			fmt.Printf(" %d. %s\n", i, e)
		}
	}

	header := metadata.New(map[string]string{"location": "MTV", "timestamp": time.Now().Format(timestampFormat)})
	stream.SendHeader(header)

	fmt.Printf("request recevied: %v\n", in)

	for i := 0; i < streamingCount; i++ {
		fmt.Printf("echo message: %v\n", in.Message)
		err := stream.Send(&api.EchoResponse{Message: in.Message})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *server) ClientStreamingEcho(stream api.Echo_ClientStreamingEchoServer) error {
	fmt.Printf("--- ClientStreamingEcho ---")
	defer func() {
		trailer := metadata.Pairs("timestamp", time.Now().Format(timestampFormat))
		stream.SetTrailer(trailer)
	}()

	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return status.Errorf(codes.DataLoss, "ClientStreamingEcho: failed to get metadata")
	}

	if t, ok := md["timestamp"]; ok {
		fmt.Printf("timestamp form metadata:\n")
		for i, e := range t {
			fmt.Printf(" %d. %s\n", i, e)
		}
	}

	header := metadata.New(map[string]string{"location": "MTV", "timestamp": time.Now().Format(timestampFormat)})
	stream.SendHeader(header)

	var message string
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			fmt.Printf("echo last received message\n")
			return stream.SendAndClose(&api.EchoResponse{Message: message})
		}
		message = in.Message
		fmt.Printf("request received: %v, building echo\n", in)
		if err != nil {
			return err
		}
	}
}

func (s *server) BidirectionalStreamingEcho(stream api.Echo_BidirectionalStreamingEchoServer) error {
	fmt.Printf("--- BidirectionalStreamingEcho ---\n")
	defer func() {
		trailer := metadata.Pairs("timestamp", time.Now().Format(timestampFormat))
		stream.SetTrailer(trailer)
	}()

	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return status.Errorf(codes.DataLoss, time.Now().Format(timestampFormat))
	}

	if t, ok := md["timestamp"]; ok {
		for i, e := range t {
			fmt.Printf(" %d. %s\n", i, e)
		}
	}

	header := metadata.New(map[string]string{"location": "MTV", "timestamp": time.Now().Format(timestampFormat)})
	stream.SendHeader(header)

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		fmt.Printf("request received %v, sending echo\n", in)
		if err := stream.Send(&api.EchoResponse{Message: in.Message}); err != nil {
			return err
		}
	}
}

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("server listening at %v\n", lis.Addr())

	s := grpc.NewServer()
	api.RegisterEchoServer(s, new(server))
	s.Serve(lis)
}
