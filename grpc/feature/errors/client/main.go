package main

import (
	"context"
	"feature/api"
	"flag"
	"log"
	"time"

	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var addr = flag.String("addr", "localhost:50051", "the address to connect to")

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func() {
		if e := conn.Close(); err != nil {
			log.Printf("failed to close connection: %s", e)
		}
	}()
	c := api.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SayHello(ctx, &api.HelloRequest{Name: "world"})
	if err != nil {
		s := status.Convert(err)

		for _, d := range s.Details() {
			switch info := d.(type) {
			case *epb.QuotaFailure:
				log.Printf("Quota failure: %s", info)
			default:
				log.Printf("Unexpected type: %s", info)
			}
		}
	}
	log.Printf("Greeting: %s", r.Message)
}
