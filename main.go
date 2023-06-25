package main

import (
	"context"
	"expvar"
	"flag"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/DazWilkin/koyeb-exporter/collector"
	"github.com/koyeb/koyeb-api-client-go/api/v1/koyeb"
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

	cfg := koyeb.NewConfiguration()
	client := koyeb.NewAPIClient(cfg)

	ctx := context.Background()
	ctx = context.WithValue(ctx, koyeb.ContextAccessToken, token)

	registry := prometheus.NewRegistry()
	registry.MustRegister(collector.NewExporterCollector(OSVersion, GoVersion, GitCommit, StartTime))

	registry.MustRegister(collector.NewAppsCollector(ctx, client))
	registry.MustRegister(collector.NewCredentialsCollector(ctx, client))
	registry.MustRegister(collector.NewDeploymentsCollector(ctx, client))
	registry.MustRegister(collector.NewDomainsCollector(ctx, client))
	registry.MustRegister(collector.NewInstancesCollector(ctx, client))
	registry.MustRegister(collector.NewSecretsCollector(ctx, client))
	registry.MustRegister(collector.NewServicesCollector(ctx, client))

	mux := http.NewServeMux()
	mux.Handle("/varz", expvar.Handler())
	mux.Handle(*metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	log.Printf("[main] Server starting (%s)", *endpoint)
	log.Printf("[main] metrics served on: %s", *metricsPath)
	log.Fatal(http.ListenAndServe(*endpoint, mux))
}
