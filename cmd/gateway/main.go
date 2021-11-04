package main

import (
	"log"

	"github.com/YaroslavGaponov/keva/internal/cluster/udp"
	"github.com/YaroslavGaponov/keva/internal/gateway"
)

func main() {

	subscriber := udp.NewSubscreber(udp.OptionsFromEnv())
	if err := subscriber.Start(); err != nil {
		panic(err)
	}

	gateway := gateway.New(subscriber)

	log.Fatal(gateway.Start())
}
