package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/setheck/smartthings-exporter/smartthings"
)

var (
	Version = "dev"
	Built   = ""
	Commit  = ""

	Banner = strings.ReplaceAll(`
                          _   _   _     _                                                  _
 ___ _ __ ___   __ _ _ __| |_| |_| |__ (_)_ __   __ _ ___        _____  ___ __   ___  _ __| |_ ___ _ __ 
/ __| '_ q _ \ / _q | '__| __| __| '_ \| | '_ \ / _q / __|_____ / _ \ \/ / '_ \ / _ \| '__| __/ _ \ '__|
\__ \ | | | | | (_| | |  | |_| |_| | | | | | | | (_| \__ \_____|  __/>  <| |_) | (_) | |  | ||  __/ |
|___/_| |_| |_|\__,_|_|   \__|\__|_| |_|_|_| |_|\__, |___/      \___/_/\_\ .__/ \___/|_|   \__\___|_|
                                                |___/                    |_|`, "q", "`")
)

type Configuration struct {
	Port     int    `envconfig:"PORT" default:"9119"`
	ApiToken string `envconfig:"API_TOKEN"`
}

var (
	ver = flag.Bool("version", false, "print version and exit")
)

func main() {
	flag.Parse()
	fmt.Println(Banner)
	fmt.Println("version:", Version)
	fmt.Println("  built:", Built)
	fmt.Println(" commit:", Commit)
	if *ver {
		os.Exit(0)
	}

	log.Println("starting up")
	var config Configuration
	if err := envconfig.Process("STE", &config); err != nil {
		log.Fatal(err)
	}

	log.Println("creating smartthings client")
	client := smartthings.NewClient(config.ApiToken, nil)
	if _, err := client.ListDevices(context.Background()); err != nil {
		log.Fatal("failed to initialize client:", err)
	}

	log.Println("creating collector")
	collector := NewCollector(client)

	prometheus.MustRegister(collector)
	http.HandleFunc("/", rootHandler)
	http.Handle("/metrics", promhttp.Handler())

	addr := fmt.Sprintf("0.0.0.0:%d", config.Port)
	log.Println("starting server on", addr)
	if err := http.ListenAndServe(addr, nil); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func rootHandler(w http.ResponseWriter, _ *http.Request) {
	// TODO:(smt) default page with info and usage etc...
	if _, err := fmt.Fprint(w, "OK"); err != nil {
		log.Println("responding to root handler request failed:", err)
	}
}
