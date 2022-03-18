package main

import (
	"context"
	"feature/api"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/alts"
)

var addr = flag.String("addr", "localhost:50051", "the address to connect to")

func callUnaryEcho(client api.EchoClient, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := client.UnaryEcho(ctx, &api.EchoRequest{Message: message})
	if err != nil {
		log.Fatalf("client.UnaryEcho(_) = _, %v: ", err)
	}
	fmt.Println("UnaryEcho: ", resp.Message)
}

func main() {
	flag.Parse()

	altsTC := alts.NewClientCreds(alts.DefaultClientOptions())

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(altsTC))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := api.NewEchoClient(conn)
	callUnaryEcho(c, "hello world")
}
