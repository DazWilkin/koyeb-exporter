package collector

import (
	"context"
	"log"

	"github.com/koyeb/koyeb-api-client-go/api/v1/koyeb"
	"github.com/prometheus/client_golang/prometheus"
)

// Ensure that ServicesCollector implements Prometheus' Collector interface
var _ prometheus.Collector = (*ServicesCollector)(nil)

// ServicesCollector collects Koyeb Services metrics
type ServicesCollector struct {
	Token string

	Up *prometheus.Desc
}

// NewServicesCollector is a function that creates a new ServicesCollector
func NewServicesCollector(token string) *ServicesCollector {
	subsystem := "services"
	return &ServicesCollector{
		Token: token,

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
	cfg := koyeb.NewConfiguration()
	client := koyeb.NewAPIClient(cfg)

	ctx := context.Background()
	ctx = context.WithValue(ctx, koyeb.ContextAccessToken, c.Token)

	rqst := client.ServicesApi.ListServices(ctx)
	resp, _, err := rqst.Execute()
	if err != nil {
		msg := "unable to list Services"
		log.Printf(msg, err)
		return
	}

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
