package collector

import (
	"context"
	"log"

	"github.com/koyeb/koyeb-api-client-go/api/v1/koyeb"
	"github.com/prometheus/client_golang/prometheus"
)

// Ensure that DeploymentCollector implements Prometheus' Collector interface
var _ prometheus.Collector = (*DomainsCollector)(nil)

// DomainsCollector collects Koyeb Domains metrics
type DomainsCollector struct {
	Ctx    context.Context
	Client *koyeb.APIClient

	Up *prometheus.Desc
}

// NewDomainsCollector is a function that creates a new DomainsCollector
func NewDomainsCollector(ctx context.Context, client *koyeb.APIClient) *DomainsCollector {
	subsystem := "domains"
	return &DomainsCollector{
		Ctx:    ctx,
		Client: client,

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
	rqst := c.Client.DomainsApi.ListDomains(c.Ctx)
	resp, _, err := rqst.Execute()
	if err != nil {
		msg := "unable to list Domains"
		log.Printf(msg, err)
		return
	}

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
