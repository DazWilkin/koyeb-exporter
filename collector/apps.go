package collector

import (
	"context"
	"log/slog"

	"github.com/DazWilkin/go-probe/probe"
	"github.com/koyeb/koyeb-api-client-go/api/v1/koyeb"
	"github.com/prometheus/client_golang/prometheus"
)

// Ensure that AppsCollector implements Prometheus' Collector interface
var _ prometheus.Collector = (*AppsCollector)(nil)

// AppsCollector collects Koyeb Apps metrics
type AppsCollector struct {
	ctx    context.Context
	client *koyeb.APIClient
	ch     chan<- probe.Status
	logger *slog.Logger

	Up *prometheus.Desc
}

// NewAppsCollector is a function that creates a new AppsCollector
func NewAppsCollector(ctx context.Context, client *koyeb.APIClient, ch chan<- probe.Status, l *slog.Logger) *AppsCollector {
	subsystem := "apps"
	logger := l.With("collector", subsystem)

	return &AppsCollector{
		ctx:    ctx,
		client: client,
		ch:     ch,
		logger: logger,

		Up: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "up"),
			"1 if the App is up, 0 otherwise",
			[]string{
				"id",
				"name",
				"organization",
				"status",
			},
			nil,
		),
	}
}

// Collect implements Prometheus' Collector interface and is used to collect metrics
func (c *AppsCollector) Collect(ch chan<- prometheus.Metric) {
	logger := c.logger.With("method", "collect")

	rqst := c.client.AppsApi.ListApps(c.ctx)
	resp, _, err := rqst.Execute()
	if err != nil {
		msg := "unable to list Apps"
		logger.Info(msg, "err", err)

		// Send probe unhealthy
		status := probe.Status{
			Healthy: false,
			Message: msg,
		}
		c.ch <- status

		return
	}

	// Send probe healthy status
	status := probe.Status{
		Healthy: true,
		Message: "ok",
	}
	c.ch <- status

	for _, app := range resp.Apps {
		ch <- prometheus.MustNewConstMetric(
			c.Up,
			prometheus.GaugeValue,
			1.0,
			[]string{
				app.GetId(),
				app.GetName(),
				app.GetOrganizationId(),
				string(app.GetStatus()),
			}...,
		)
	}

}

// Describe implements Prometheus' Collector interface and is used to describe metrics
func (c *AppsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Up
}
