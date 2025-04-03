package collector

import (
	"context"
	"log/slog"

	"github.com/DazWilkin/go-probe/probe"
	"github.com/koyeb/koyeb-api-client-go/api/v1/koyeb"
	"github.com/prometheus/client_golang/prometheus"
)

// Ensure that DeploymentCollector implements Prometheus' Collector interface
var _ prometheus.Collector = (*DomainsCollector)(nil)

// DomainsCollector collects Koyeb Domains metrics
type DomainsCollector struct {
	ctx    context.Context
	client *koyeb.APIClient
	ch     chan<- probe.Status
	logger *slog.Logger

	Up *prometheus.Desc
}

// NewDomainsCollector is a function that creates a new DomainsCollector
func NewDomainsCollector(ctx context.Context, client *koyeb.APIClient, ch chan<- probe.Status, l *slog.Logger) *DomainsCollector {
	subsystem := "domains"
	logger := l.With("collector", subsystem)

	return &DomainsCollector{
		ctx:    ctx,
		client: client,
		ch:     ch,
		logger: logger,

		Up: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "up"),
			"1 if the Domain is up, 0 otherwise",
			[]string{
				"id",
				"app_id",
				"organization_id",
				"name",
				"status",
				"type",
			},
			nil,
		),
	}
}

// Collect implements Prometheus' Collector interface and is used to collect metrics
func (c *DomainsCollector) Collect(ch chan<- prometheus.Metric) {
	logger := c.logger.With("method", "collect")

	rqst := c.client.DomainsApi.ListDomains(c.ctx)
	resp, _, err := rqst.Execute()
	if err != nil {
		msg := "unable to list Domains"
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

	for _, domain := range resp.Domains {
		ch <- prometheus.MustNewConstMetric(
			c.Up,
			prometheus.GaugeValue,
			1.0,
			[]string{
				domain.GetId(),
				domain.GetAppId(),
				domain.GetOrganizationId(),
				domain.GetName(),
				string(domain.GetStatus()),
				string(domain.GetType()),
			}...,
		)
	}

}

// Describe implements Prometheus' Collector interface and is used to describe metrics
func (c *DomainsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Up
}
