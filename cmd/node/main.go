package main

import (
	"log"

	"github.com/YaroslavGaponov/keva/internal/cluster/udp"
	"github.com/YaroslavGaponov/keva/internal/node"
	"github.com/YaroslavGaponov/keva/internal/storage/memory"
)

func main() {

	publisher := udp.NewPublisher(udp.OptionsFromEnv())

	storage := memory.New(memory.OptionsFromEnv())
	if err := storage.Open(); err != nil {
		panic(err)
	}

	node, err := node.New(publisher, storage)
	if err != nil {
		panic(err)
	}

	log.Fatal(node.Start())
}
