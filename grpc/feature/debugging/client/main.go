package main

import (
	"context"
	"feature/api"
	"flag"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/channelz/service"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/resolver/manual"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()

	/***** Set up the server serving channelz service. *****/
	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	s := grpc.NewServer()
	service.RegisterChannelzServiceToServer(s)
	go s.Serve(lis)
	defer s.Stop()

	/***** Initialize manual resolver and Dial *****/
	r := manual.NewBuilderWithScheme("whatever")
	// Set up a connection to the server.
	conn, err := grpc.Dial(r.Scheme()+":///test.server",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithResolvers(r),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	r.UpdateState(resolver.State{
		Addresses: []resolver.Address{
			{Addr: ":10001"},
			{Addr: ":10002"},
			{Addr: ":10003"},
		},
	})

	c := api.NewGreeterClient(conn)

	for i := 0; i < 100; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
		defer cancel()

		r, err := c.SayHello(ctx, &api.HelloRequest{Name: *name})
		if err != nil {
			log.Printf("could not greet: %v", err)
		} else {
			log.Printf("Greeting: %s", r.Message)
		}
	}

	select {}
}
