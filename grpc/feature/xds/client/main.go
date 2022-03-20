package main

import (
	"context"
	"feature/api"
	"flag"
	"log"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	xdscreds "google.golang.org/grpc/credentials/xds"

	_ "google.golang.org/grpc/xds"
)

var (
	target   = flag.String("target", "xds:///localhost:50051", "uri of the Greeter Server, e.g. 'xds:///helloworld-service:8080'")
	name     = flag.String("name", "world", "name you wished to be greeted by the server")
	xdsCreds = flag.Bool("xds_creds", false, "whether the server should use xDS APIs to receive security configuration")
)

func main() {
	flag.Parse()

	if !strings.HasPrefix(*target, "xds:///") {
		log.Fatalf("-target must use a URI with scheme set to 'xds'")
	}

	creds := insecure.NewCredentials()
	if *xdsCreds {
		log.Println("Using xDS credentials...")
		var err error
		if creds, err = xdscreds.NewClientCredentials(xdscreds.ClientOptions{
			FallbackCreds: insecure.NewCredentials(),
		}); err != nil {
			log.Fatalf("failed to create client-side xDS credentials: %v", err)
		}
	}

	conn, err := grpc.Dial(*target, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("grpc.Dial(%s) failed: %v", *target, err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	c := api.NewGreeterClient(conn)
	r, err := c.SayHello(ctx, &api.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
		log.Printf("Greeting: %s", r.GetMessage())
	}
}
