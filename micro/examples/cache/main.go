package main

import (
	"examples/cache/api"
	"examples/cache/handler"
	"log"

	"go-micro.dev/v4"
)

var (
	service = "go.micro.srv.cache"
	version = "latest"
)

func main() {
	srv := micro.NewService(
		micro.Name(service),
		micro.Version(version),
	)

	srv.Init()

	api.RegisterCacheHandler(srv.Server(), handler.NewCache())

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
