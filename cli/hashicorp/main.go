package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mitchellh/cli"
)

func main() {
	c := cli.NewCLI("app", "1.0.0")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"hello": helloCommandFactory,
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}

func helloCommandFactory() (cli.Command, error) {
	return helloCommand{}, nil
}

type helloCommand struct {
}

func (helloCommand) Help() string {
	return "hello world"
}

func (helloCommand) Run(args []string) int {
	fmt.Fprintln(os.Stdout, "hello", args[0])
	return len(args)
}

func (helloCommand) Synopsis() string {
	return "this is a helloworld command"
}
