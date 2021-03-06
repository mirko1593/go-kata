package main

import (
	"fmt"
	"log"
	"time"

	"go-micro.dev/v4/broker"
	"go-micro.dev/v4/cmd"
)

var (
	topic = "go.micro.topic.foo"
)

func pub() {
	tick := time.NewTicker(time.Second)
	i := 0

	for _ = range tick.C {
		msg := &broker.Message{
            Header: map[string]string{
                "id": fmt.Sprintf("%d", i)
            },
            Body: []byte(fmt.Sprintf("%d: %s", i, time.Now().String())),
        }
        i++
	}
}

func main() {
    cmd.Init()

    if err := broker.Init(); err != nil {
        log.Fatalf("Broker Init error: %v", err)
    }

    if err := broker.Connect(); err != nil {
        log.Fatalf("Broker Connect error:%v", err)
    }

    pub()
}
