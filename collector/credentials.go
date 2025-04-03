package collector

import (
	"context"
	"log/slog"

	"github.com/DazWilkin/go-probe/probe"
	"github.com/koyeb/koyeb-api-client-go/api/v1/koyeb"
	"github.com/prometheus/client_golang/prometheus"
)

// Ensure that DeploymentCollector implements Prometheus' Collector interface
var _ prometheus.Collector = (*CredentialsCollector)(nil)

// CredentialsCollector collects Koyeb Credentials metrics
type CredentialsCollector struct {
	ctx    context.Context
	client *koyeb.APIClient
	ch     chan<- probe.Status
	logger *slog.Logger

	Up *prometheus.Desc
}

// NewCredentialsCollector is a function that creates a new CredentialsCollector
func NewCredentialsCollector(ctx context.Context, client *koyeb.APIClient, ch chan<- probe.Status, l *slog.Logger) *CredentialsCollector {
	subsystem := "credentials"
	logger := l.With("collector", subsystem)

	return &CredentialsCollector{
		ctx:    ctx,
		client: client,
		ch:     ch,
		logger: logger,

		Up: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "up"),
			"1 if the Credentials is up, 0 otherwise",
			[]string{
				"id",
				"organization_id",
				"user_id",
				"name",
			},
			nil,
		),
	}
}

// Collect implements Prometheus' Collector interface and is used to collect metrics
func (c *CredentialsCollector) Collect(ch chan<- prometheus.Metric) {
	logger := c.logger.With("method", "collect")

	rqst := c.client.CredentialsApi.ListCredentials(c.ctx)
	resp, _, err := rqst.Execute()
	if err != nil {
		msg := "unable to list Credentials"
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

	for _, credential := range resp.Credentials {
		ch <- prometheus.MustNewConstMetric(
			c.Up,
			prometheus.GaugeValue,
			1.0,
			[]string{
				credential.GetId(),
				credential.GetOrganizationId(),
				credential.GetUserId(),
				credential.GetName(),
			}...,
		)
	}

}

// Describe implements Prometheus' Collector interface and is used to describe metrics
func (c *CredentialsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Up
}
