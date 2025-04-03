package collector

import (
	"context"
	"log/slog"

	"github.com/DazWilkin/go-probe/probe"
	"github.com/koyeb/koyeb-api-client-go/api/v1/koyeb"
	"github.com/prometheus/client_golang/prometheus"
)

// Ensure that ServicesCollector implements Prometheus' Collector interface
var _ prometheus.Collector = (*ServicesCollector)(nil)

// ServicesCollector collects Koyeb Services metrics
type ServicesCollector struct {
	ctx    context.Context
	client *koyeb.APIClient
	ch     chan<- probe.Status
	logger *slog.Logger

	Up *prometheus.Desc
}

// NewServicesCollector is a function that creates a new ServicesCollector
func NewServicesCollector(ctx context.Context, client *koyeb.APIClient, ch chan<- probe.Status, l *slog.Logger) *ServicesCollector {
	subsystem := "services"
	logger := l.With("collector", subsystem)

	return &ServicesCollector{
		ctx:    ctx,
		client: client,
		ch:     ch,
		logger: logger,

		Up: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "up"),
			"1 if the App is up, 0 otherwise",
			[]string{
				"id",
				"app_id",
				"organization_id",
				"name",
				"status",
			},
			nil,
		),
	}
}

// Collect implements Prometheus' Collector interface and is used to collect metrics
func (c *ServicesCollector) Collect(ch chan<- prometheus.Metric) {
	logger := c.logger.With("method", "collect")

	rqst := c.client.ServicesApi.ListServices(c.ctx)
	resp, _, err := rqst.Execute()
	if err != nil {
		msg := "unable to list Services"
		logger.Error(msg, "err", err)

		// Send probe unhealthy status
		// Doesn't surface the API error message (should it!?)
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

	for _, service := range resp.Services {
		ch <- prometheus.MustNewConstMetric(
			c.Up,
			prometheus.GaugeValue,
			1.0,
			[]string{
				service.GetId(),
				service.GetAppId(),
				service.GetOrganizationId(),
				service.GetName(),
				string(service.GetStatus()),
			}...,
		)
	}

}

// Describe implements Prometheus' Collector interface and is used to describe metrics
func (c *ServicesCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Up
}
