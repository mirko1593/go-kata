package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	yaml "github.com/asim/go-micro/plugins/config/encoder/yaml/v4"
	proto "github.com/asim/go-micro/plugins/config/source/grpc/v4/proto"
	"go-micro.dev/v4/config"
	"go-micro.dev/v4/config/reader"
	"go-micro.dev/v4/config/reader/json"
	"go-micro.dev/v4/config/source/file"
	"go-micro.dev/v4/logger"
	"google.golang.org/grpc"
)

var (
	mux        sync.RWMutex
	configMaps = make(map[string]*proto.ChangeSet)
	apps       = []string{"micro", "extra"}
	cfg        config.Config
)

// Service ...
type Service struct{}

func main() {
	enc := yaml.NewEncoder()
	cfg, _ = config.NewConfig(config.WithReader(json.NewReader(
		reader.WithEncoder(enc),
	)))

	err := loadConfigFile()

	server := grpc.NewServer()
	proto.RegisterSourceServer(server, new(Service))
	ts, err := net.Listen("tcp", ":8600")
	if err != nil {
		logger.Fatal(err)
	}

	logger.Infof("configServer started")
	err = server.Serve(ts)
	if err != nil {
		logger.Fatal(err)
	}
}

func loadConfigFile() (err error) {
	for _, app := range apps {
		if err := cfg.Load(file.NewSource(
			file.WithPath("./conf/" + app + ".yaml"),
		)); err != nil {
			log.Fatalf("[loadConfigFile] load files error: %v", err)
			return err
		}
	}

	watcher, err := cfg.Watch()
	if err != nil {
		log.Fatalf("[loadConfigFile] start watching files error, %s", err)
		return err
	}

	go func() {
		for {
			v, err := watcher.Next()
			if err != nil {
				log.Fatalf("[loadConfigFile] watch files erorr: %s", err)
				return
			}

			logger.Infof("[loadConfigFile] file change, %s", string(v.Bytes()))
		}
	}()

	return
}

func (s Service) Read(ctx context.Context, req *proto.ReadRequest) (rsp *proto.ReadResponse, err error) {
	appName := parsePath(req.Path)
	switch appName {
	case "micro", "extra":
		rsp = &proto.ReadResponse{
			ChangeSet: getConfig(appName),
		}
		return
	default:
		err = fmt.Errorf("[Read] the first path is invalid")
		return
	}
}

func (s Service) Watch(req *proto.WatchRequest, server proto.Source_WatchServer) (err error) {
	appName := parsePath(req.Path)
	rsp := &proto.WatchResponse{
		ChangeSet: getConfig(appName),
	}

	if err := server.Send(rsp); err != nil {
		logger.Infof("[Watch] watch files error, %s", err)
		return err
	}

	return
}

func getConfig(appName string) *proto.ChangeSet {
	bytes := cfg.Get(appName).Bytes()

	logger.Infof("[getConfig] appName %s", string(bytes))

	return &proto.ChangeSet{
		Data:      bytes,
		Checksum:  fmt.Sprintf("%x", md5.Sum(bytes)),
		Format:    "yml",
		Source:    "file",
		Timestamp: time.Now().Unix(),
	}

	return nil
}

func parsePath(path string) (appName string) {
	paths := strings.Split(path, "/")

	if paths[0] == "" && len(paths) > 1 {
		return paths[1]
	}

	return paths[0]
}
