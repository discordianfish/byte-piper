package main

import (
	"flag"
	"log"

	"github.com/docker-infra/bytes-piper/pipeline"
)

var (
	config = flag.String("c", "", "path to config")
)

func main() {
	flag.Parse()
	if *config == "" {
		log.Fatal("No config provided")
	}
	pipe, err := pipeline.New(*config)
	if err != nil {
		log.Fatal(err)
	}

	if err := pipe.Run(); err != nil {
		log.Fatal(err)
	}
}
