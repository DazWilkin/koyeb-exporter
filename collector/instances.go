package collector

import (
	"context"
	"log/slog"

	"github.com/DazWilkin/go-probe/probe"
	"github.com/koyeb/koyeb-api-client-go/api/v1/koyeb"
	"github.com/prometheus/client_golang/prometheus"
)

// Ensure that InstancesCollector implements Prometheus' Collector interface
var _ prometheus.Collector = (*InstancesCollector)(nil)

// InstancesCollector collects Koyeb Apps metrics
type InstancesCollector struct {
	ctx    context.Context
	client *koyeb.APIClient
	ch     chan<- probe.Status
	logger *slog.Logger

	Up *prometheus.Desc
}

// NewInstancesCollector is a function that creates a new InstancesCollector
func NewInstancesCollector(ctx context.Context, client *koyeb.APIClient, ch chan<- probe.Status, l *slog.Logger) *InstancesCollector {
	subsystem := "instances"
	logger := l.With("collector", subsystem)

	return &InstancesCollector{
		ctx:    ctx,
		client: client,
		ch:     ch,
		logger: logger,

		Up: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "up"),
			"1 if the Instance is up, 0 otherwise",
			[]string{
				"id",
				"app_id",
				"service_id",
				"organization_id",
				"region",
				"status",
			},
			nil,
		),
	}
}

// Collect implements Prometheus' Collector interface and is used to collect metrics
func (c *InstancesCollector) Collect(ch chan<- prometheus.Metric) {
	logger := c.logger.With("method", "collect")

	rqst := c.client.InstancesApi.ListInstances(c.ctx)
	resp, _, err := rqst.Execute()
	if err != nil {
		msg := "unable to list Instances"
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

	for _, instance := range resp.Instances {
		ch <- prometheus.MustNewConstMetric(
			c.Up,
			prometheus.GaugeValue,
			1.0,
			[]string{
				instance.GetId(),
				instance.GetAppId(),
				instance.GetServiceId(),
				instance.GetOrganizationId(),
				instance.GetRegion(),
				string(instance.GetStatus()),
			}...,
		)
	}

}

// Describe implements Prometheus' Collector interface and is used to describe metrics
func (c *InstancesCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Up
}
