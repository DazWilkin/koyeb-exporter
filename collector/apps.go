package collector

import (
	"context"
	"log"

	"github.com/koyeb/koyeb-api-client-go/api/v1/koyeb"
	"github.com/prometheus/client_golang/prometheus"
)

// Ensure that AppsCollector implements Prometheus' Collector interface
var _ prometheus.Collector = (*AppsCollector)(nil)

// AppsCollector collects Koyeb Apps metrics
type AppsCollector struct {
	Ctx    context.Context
	Client *koyeb.APIClient

	Up *prometheus.Desc
}

// NewAppsCollector is a function that creates a new AppsCollector
func NewAppsCollector(ctx context.Context, client *koyeb.APIClient) *AppsCollector {
	subsystem := "apps"
	return &AppsCollector{
		Ctx:    ctx,
		Client: client,

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
	rqst := c.Client.AppsApi.ListApps(c.Ctx)
	resp, _, err := rqst.Execute()
	if err != nil {
		msg := "unable to list Apps"
		log.Printf(msg, err)
		return
	}

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
