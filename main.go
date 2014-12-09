package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/docker-infra/bytes-piper/pipeline"
)

var (
	debugEndpoint = flag.String("d", "", "Enable pprof debugging endpoint on given host:port")
	loop          = flag.Duration("l", 0, "Repeat pipelines at given interval")
	plines        pipelines
)

type pipelines []*pipeline.Pipeline

func (p *pipelines) String() string {
	return fmt.Sprintf("%v", *p)
}
func (p *pipelines) Set(v string) error {
	pipe, err := pipeline.New(v)
	if err != nil {
		return err
	}
	*p = append(*p, pipe)
	return nil
}

func main() {
	var listenErr chan error
	if *debugEndpoint != "" {
		go func() {
			listenErr <- http.ListenAndServe(*debugEndpoint, nil)

		}()
	}
	flag.Var(&plines, "c", "Path to config, may be repeated")
	flag.Parse()
	if len(plines) == 0 {
		log.Fatal("No configs provided")
	}

	for {
		for _, p := range plines {
			if err := p.Run(); err != nil {
				log.Fatal("Error running pipeline: ", err)
			}
		}
		if *loop == 0 {
			break
		}
		plines = pipelines{}
		log.Print("Sleeping for ", *loop)
		time.Sleep(*loop)
		flag.Parse() // We can do that nicer..
	}
	if *debugEndpoint != "" {
		log.Print("Debugging enabled, keep listening for debugging")
		log.Print(<-listenErr)
	}
}
