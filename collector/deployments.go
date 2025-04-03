package collector

import (
	"context"
	"log/slog"

	"github.com/DazWilkin/go-probe/probe"
	"github.com/koyeb/koyeb-api-client-go/api/v1/koyeb"
	"github.com/prometheus/client_golang/prometheus"
)

// Ensure that DeploymentCollector implements Prometheus' Collector interface
var _ prometheus.Collector = (*DeploymentsCollector)(nil)

// DeploymentsCollector collects Koyeb Deployments metrics
type DeploymentsCollector struct {
	ctx    context.Context
	client *koyeb.APIClient
	ch     chan<- probe.Status
	logger *slog.Logger

	Up *prometheus.Desc
}

// NewDeploymentsCollector is a function that creates a new DeploymentsCollector
func NewDeploymentsCollector(ctx context.Context, client *koyeb.APIClient, ch chan<- probe.Status, l *slog.Logger) *DeploymentsCollector {
	subsystem := "deployments"
	logger := l.With("collector", subsystem)

	return &DeploymentsCollector{
		ctx:    ctx,
		client: client,
		ch:     ch,
		logger: logger,

		Up: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "up"),
			"1 if the Deployment is up, 0 otherwise",
			[]string{
				"id",
				"app_id",
				"deployment_group",
				"name",
				"service_id",
				"status",
				"type",
			},
			nil,
		),
	}
}

// Collect implements Prometheus' Collector interface and is used to collect metrics
func (c *DeploymentsCollector) Collect(ch chan<- prometheus.Metric) {
	logger := c.logger.With("method", "collect")

	rqst := c.client.DeploymentsApi.ListDeployments(c.ctx)
	resp, _, err := rqst.Execute()
	if err != nil {
		msg := "unable to list Deployments"
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

	for _, deployment := range resp.Deployments {
		ch <- prometheus.MustNewConstMetric(
			c.Up,
			prometheus.GaugeValue,
			1.0,
			[]string{
				deployment.GetId(),
				deployment.GetAppId(),
				deployment.GetDeploymentGroup(),
				deployment.Definition.GetName(),
				deployment.GetServiceId(),
				string(deployment.GetStatus()),
				string(deployment.Definition.GetType()),
			}...,
		)
	}

}

// Describe implements Prometheus' Collector interface and is used to describe metrics
func (c *DeploymentsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Up
}
