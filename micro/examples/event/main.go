package main

import (
	"context"
	"examples/event/api"

	"go-micro.dev/v4"
	"go-micro.dev/v4/util/log"
)

// Event ...
type Event struct{}

// Process ...
func (e *Event) Process(ctx context.Context, event *api.Event) error {
	log.Logf("Received event %+v\n", event)
	return nil
}

func main() {
	srv := micro.NewService(
		micro.Name("user"),
	)
	srv.Init()

	micro.RegisterSubscriber("go.micro.evt.user", srv.Server(), new(Event))

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
