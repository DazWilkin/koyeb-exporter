package collector

import (
	"context"
	"log/slog"

	"github.com/DazWilkin/go-probe/probe"
	"github.com/DazWilkin/koyeb-exporter/types"
	"github.com/koyeb/koyeb-api-client-go/api/v1/koyeb"
	"github.com/prometheus/client_golang/prometheus"
)

// Ensure that DeploymentCollector implements Prometheus' Collector interface
var _ prometheus.Collector = (*SecretsCollector)(nil)

// SecretsCollector collects Koyeb Secrets metrics
type SecretsCollector struct {
	ctx    context.Context
	client *koyeb.APIClient
	ch     chan<- probe.Status
	logger *slog.Logger

	Up *prometheus.Desc
}

// NewSecretsCollector is a function that creates a new SecretsCollector
func NewSecretsCollector(ctx context.Context, client *koyeb.APIClient, ch chan<- probe.Status, l *slog.Logger) *SecretsCollector {
	subsystem := "secrets"
	logger := l.With("collector", subsystem)

	return &SecretsCollector{
		ctx:    ctx,
		client: client,
		ch:     ch,
		logger: logger,

		Up: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "up"),
			"1 if the Secret is up, 0 otherwise",
			[]string{
				"id",
				"organization_id",
				"name",
				"type",
				"registry", // Synthetic
			},
			nil,
		),
	}
}

// Collect implements Prometheus' Collector interface and is used to collect metrics
func (c *SecretsCollector) Collect(ch chan<- prometheus.Metric) {
	logger := c.logger.With("method", "collect")

	rqst := c.client.SecretsApi.ListSecrets(c.ctx)
	resp, _, err := rqst.Execute()
	if err != nil {
		msg := "unable to list Secrets"
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

	for _, secret := range resp.Secrets {
		ch <- prometheus.MustNewConstMetric(
			c.Up,
			prometheus.GaugeValue,
			1.0,
			[]string{
				secret.GetId(),
				secret.GetOrganizationId(),
				secret.GetName(),
				string(secret.GetType()),
				types.GetRegistryType(secret).String(),
			}...,
		)
	}

}

// Describe implements Prometheus' Collector interface and is used to describe metrics
func (c *SecretsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Up
}
