package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/setheck/smartthings-exporter/smartthings"
)

var (
	Version = "dev"
	Built   = ""
	Commit  = ""
)

type Configuration struct {
	Port     int    `envconfig:"PORT" default:"9119"`
	ApiToken string `envconfig:"API_TOKEN"`
}

func main() {
	ver := flag.Bool("version", false, "print version and exit")
	flag.Parse()
	fmt.Println("version:", Version, "built:", Built, "commit:", Commit)
	if *ver {
		os.Exit(0)
	}
	log.Println("starting up")
	var config Configuration
	if err := envconfig.Process("STE", &config); err != nil {
		log.Fatal(err)
	}

	log.Println("creating client")
	client := smartthings.NewClient(config.ApiToken)
	if _, err := client.ListDevices(); err != nil {
		log.Fatal("failed to initialize client:", err)
	}

	log.Println("creating collector")
	collector := NewCollector(client)

	prometheus.MustRegister(collector)
	http.Handle("/metrics", promhttp.Handler())

	addr := fmt.Sprintf("0.0.0.0:%d", config.Port)
	log.Println("starting server on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
