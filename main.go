package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/docker-infra/byte-piper/pipeline"
)

var (
	debugEndpoint  = flag.String("d", "", "Enable pprof debugging endpoint on given host:port")
	loop           = flag.Duration("r", 0, "Daemon mode; Repeat pipelines at given interval")
	prometheusAddr = flag.String("l", "", "Expose prometheus metrics on given host:port, requires daemon mode")
	backupDuration = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "bytes_piper_backup_duration_seconds",
		Help: "Duration of given backup pipeline",
	}, []string{"name"})
	backupSeen = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "bytes_piper_backup_last_successful",
		Help: "Last time the given backup pipeline ran successfully",
	}, []string{"name"})
	backupSize = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "bytes_piper_backup_size_bytes",
		Help: "Bytes ran through the backup pipeline",
	}, []string{"name"})
	backupsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "bytes_piper_backups_total",
		Help: "Total number of backups pipeline runs",
	}, []string{"name"})
	backupsFailed = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "bytes_piper_backups_failed_total",
		Help: "Total number of failed backups pipeline runs",
	}, []string{"name"})
	plines pipelines
)

type pipelines []string

func (p *pipelines) String() string {
	return fmt.Sprintf("%v", *p)
}
func (p *pipelines) Set(v string) error {
	*p = append(*p, v)
	return nil
}

func init() {
	prometheus.MustRegister(backupDuration)
	prometheus.MustRegister(backupSeen)
	prometheus.MustRegister(backupSize)
	prometheus.MustRegister(backupsTotal)
	prometheus.MustRegister(backupsFailed)
}

func main() {
	var listenErr chan error
	flag.Var(&plines, "c", "Path to config, may be repeated")
	flag.Parse()

	if len(plines) == 0 {
		log.Fatal("No configs provided")
	}

	if *prometheusAddr != "" {
		if *loop == 0 {
			log.Fatal("Can only expose metrics in daemon mode")
		}
		http.Handle("/metrics", prometheus.Handler())
		go http.ListenAndServe(*prometheusAddr, nil)
	}
	if *debugEndpoint != "" {
		go func() {
			listenErr <- http.ListenAndServe(*debugEndpoint, nil)

		}()
	}

	for {
		for _, file := range plines {
			log.Print("# Running ", file)
			backupsTotal.WithLabelValues(file).Inc()
			pipe, err := pipeline.New(file)
			if err != nil {
				log.Printf("ERROR loading %s: %s", file, err)
				backupsFailed.WithLabelValues(file).Inc()
				continue
			}
			begin := time.Now()
			bytesWritten, err := pipe.Run()
			if err != nil {
				log.Printf("ERROR running %s: %s", file, err)
				backupsFailed.WithLabelValues(file).Inc()
				continue
			}
			backupSize.WithLabelValues(file).Set(float64(bytesWritten))

			now := time.Now()
			backupSeen.WithLabelValues(file).Set(float64(now.Unix()))
			backupDuration.WithLabelValues(file).Set(now.Sub(begin).Seconds())
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
