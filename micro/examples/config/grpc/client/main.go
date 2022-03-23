package main

import (
	"log"

	grpcConfig "github.com/asim/go-micro/plugins/config/source/grpc/v4"
	"go-micro.dev/v4/config"
	"go-micro.dev/v4/logger"
)

// Micro ...
type Micro struct {
	Info
}

// Info ...
type Info struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
	Message string `json:"message,omitempty"`
	Age     int    `json:"age,omitempty"`
}

func main() {
	source := grpcConfig.NewSource(
		grpcConfig.WithAddress("127.0.0.1:8600"),
		grpcConfig.WithPath("/micro"),
	)

	conf, _ := config.NewConfig()

	if err := conf.Load(source); err != nil {
		log.Fatal(err)
	}

	configs := &Micro{}
	if err := conf.Scan(configs); err != nil {
		logger.Fatal(err)
	}

	logger.Infof("Read config: %s", string(conf.Bytes()))

	watcher, err := conf.Watch()
	if err != nil {
		logger.Fatal(err)
	}

	logger.Infof("Watching for changes...")

	for {
		v, err := watcher.Next()
		if err != nil {
			logger.Fatal(err)
		}

		logger.Infof("Watching for changes: %v", string(v.Bytes()))
	}

}
