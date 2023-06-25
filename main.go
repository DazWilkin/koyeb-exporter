package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/DazWilkin/koyeb-exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// GitCommit is the git commit value and is expected to be set during build
	GitCommit string
	// GoVersion is the Golang runtime version
	GoVersion = runtime.Version()
	// OSVersion is the OS version (uname --kernel-release) and is expected to be set during build
	OSVersion string
	// StartTime is the start time of the exporter represented as a UNIX epoch
	StartTime = time.Now().Unix()
)
var (
	endpoint    = flag.String("endpoint", ":8080", "The endpoint of the Expoter's HTTP server")
	metricsPath = flag.String("path", "/metrics", "The path on which Prometheus metrics will be served")
)

func healthz(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "ok")
}
func main() {
	flag.Parse()

	if GitCommit == "" {
		log.Println("[main] GitCommit value unchanged: expected to be set during build")
	}
	if OSVersion == "" {
		log.Println("[main] OSVersion value unchanged: expected to be set during build")
	}

	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("unable to get TOKEN")
	}

	registry := prometheus.NewRegistry()
	registry.MustRegister(collector.NewExporterCollector(OSVersion, GoVersion, GitCommit, StartTime))

	registry.MustRegister(collector.NewAppsCollector(token))
	registry.MustRegister(collector.NewCredentialsCollector(token))
	registry.MustRegister(collector.NewDeploymentsCollector(token))
	registry.MustRegister(collector.NewDomainsCollector(token))
	registry.MustRegister(collector.NewInstancesCollector(token))
	registry.MustRegister(collector.NewSecretsCollector(token))
	registry.MustRegister(collector.NewServicesCollector(token))

	mux := http.NewServeMux()
	mux.Handle("/healthz", http.HandlerFunc(healthz))
	mux.Handle(*metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	log.Printf("[main] Server starting (%s)", *endpoint)
	log.Printf("[main] metrics served on: %s", *metricsPath)
	log.Fatal(http.ListenAndServe(*endpoint, mux))
}
