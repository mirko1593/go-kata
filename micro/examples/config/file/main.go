package main

import (
	"fmt"

	"go-micro.dev/v4/config"
	"go-micro.dev/v4/config/source/file"
)

func main() {
	if err := config.Load(file.NewSource(
		file.WithPath("./config.json"),
	)); err != nil {
		fmt.Println(err)
		return
	}

	type Host struct {
		Address string `json:"address"`
		Port    int    `json:"port"`
	}

	var host Host

	if err := config.Get("hosts", "database").Scan(&host); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(host.Address, host.Port)
}
