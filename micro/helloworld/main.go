package main

import (
	"helloworld/handler"
	pb "helloworld/proto"

	"go-micro.dev/v4"
	log "go-micro.dev/v4/logger"
)

var (
	service = "helloworld"
	version = "latest"
)

func main() {
	// Create service
	srv := micro.NewService(
		micro.Name(service),
		micro.Version(version),
	)
	srv.Init()

	// Register handler
	pb.RegisterHelloworldHandler(srv.Server(), new(handler.Helloworld))

	// Run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
