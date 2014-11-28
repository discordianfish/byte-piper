package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/docker-infra/bytes-piper/pipeline"
)

var (
	config        = flag.String("c", "", "path to config")
	debugEndpoint = flag.String("d", "", "enable pprof debugging endpoint on given host:port")
)

func main() {
	var listenErr chan error
	if *debugEndpoint != "" {
		go func() {
			listenErr <- http.ListenAndServe(*debugEndpoint, nil)

		}()
	}

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

	if *debugEndpoint != "" {
		log.Print("Debugging enabled, keep listening for debugging")
		log.Print(<-listenErr)
	}
}
