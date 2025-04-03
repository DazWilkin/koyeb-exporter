package main

import (
	"context"
	"expvar"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/DazWilkin/go-probe/probe"
	"github.com/DazWilkin/koyeb-exporter/collector"
	"github.com/koyeb/koyeb-api-client-go/api/v1/koyeb"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	robotsTemplate string = `
User-agent: *
Disallow: /
`
	rootTemplate string = `
{{- define "content" }}
<!DOCTYPE html>
<html lang="en-US">
<head>
	<meta name="description" content="Prometheus Exporter for {{ .Name }}">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>Prometheus Exporter for {{ .Name }}</title>
	<style>
	body { font-family: Verdana; }
	.footer { margin-top: 2rem; font-size: 0.8rem; color: #6a737d; border-top: 1px solid #eaecef; padding-top: 1rem; }
	</style>
</head>
<body>
	<h2>Prometheus Exporter for {{ .Name }}</h2>
	<hr/>
	<ul>
	<li><a href="{{ .MetricsPath }}">metrics</a></li>
	<li><a href="/healthz">healthz</a></li>
	<li><a href="/varz">varz</a></li>

	</ul>
	<div class="footer">
		<p>Version: {{ .GitCommit }} | Go: {{ .GoVersion }} | OS: {{ .OSVersion }}</p>
		<p>Started: {{ .StartTimeFormatted }}</p>
	</div>
</body>
</html>
{{- end}}
`
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

type Content struct {
	Name               string
	MetricsPath        string
	GitCommit          string
	GoVersion          string
	OSVersion          string
	StartTimeFormatted string
}

func robots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(robotsTemplate)); err != nil {
		slog.Error("unable to write robots handler content")
	}
}
func root(content Content) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		t := template.Must(template.New("content").Parse(rootTemplate))
		if err := t.ExecuteTemplate(w, "content", content); err != nil {
			slog.Error("unable to execute template")
		}
	}
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	flag.Parse()

	if GitCommit == "" {
		logger.Error("value unchanged: expected GitCommit to be set during build")
	}
	if OSVersion == "" {
		logger.Error("value unchanged: expected OSVersion to be set during build")
	}

	token := os.Getenv("TOKEN")
	if token == "" {
		logger.Error("unable to get TOKEN from environment")
		return
	}

	p := probe.New("liveness", logger)
	healthz := p.Handler(logger)

	// Channel is shared by the Updater (subscriber) and the Collectors (publisher)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Start the Updater
	ch := make(chan probe.Status)
	go p.Updater(ctx, ch, nil)

	cfg := koyeb.NewConfiguration()
	client := koyeb.NewAPIClient(cfg)
	ctx = context.WithValue(ctx, koyeb.ContextAccessToken, token)

	registry := prometheus.NewRegistry()

	// Register collectors
	for _, c := range []struct {
		name      string
		collector prometheus.Collector
	}{
		{
			"exporter",
			collector.NewExporterCollector(OSVersion, GoVersion, GitCommit, StartTime),
		},
		{
			"apps",
			collector.NewAppsCollector(ctx, client, ch, logger),
		},
		{
			"credentials",
			collector.NewCredentialsCollector(ctx, client, ch, logger),
		},
		{
			"deployments",
			collector.NewDeploymentsCollector(ctx, client, ch, logger),
		},
		{
			"domains",
			collector.NewDomainsCollector(ctx, client, ch, logger),
		},
		{
			"instances",
			collector.NewInstancesCollector(ctx, client, ch, logger),
		},
		{
			"secrets",
			collector.NewSecretsCollector(ctx, client, ch, logger),
		},
		{
			"services",
			collector.NewServicesCollector(ctx, client, ch, logger),
		},
	} {
		if err := registry.Register(c.collector); err != nil {
			logger.Error("failed to register collector",
				"collector", c.name,
				"err", err,
			)
		}
	}

	mux := http.NewServeMux()

	// Create content for the root page
	content := Content{
		Name:               "Koyeb",
		MetricsPath:        *metricsPath,
		GitCommit:          GitCommit,
		GoVersion:          GoVersion,
		OSVersion:          OSVersion,
		StartTimeFormatted: time.Unix(StartTime, 0).Format(time.RFC3339),
	}

	mux.Handle("/", root(content))
	mux.Handle("/healthz", http.HandlerFunc(healthz))
	mux.Handle("/robots.txt", http.HandlerFunc(robots))

	mux.Handle("/varz", expvar.Handler())
	mux.Handle(*metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	logger.Info("Server starting",
		"endpoint", *endpoint,
		"metrics", *metricsPath,
	)

	server := &http.Server{
		Addr:         *endpoint,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	logger.Error("unable to server",
		"err", server.ListenAndServe(),
	)
}
