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
	Token string

	Up *prometheus.Desc
}

// NewDomainsCollector is a function that creates a new DomainsCollector
func NewDomainsCollector(token string) *DomainsCollector {
	subsystem := "domains"
	return &DomainsCollector{
		Token: token,

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
	cfg := koyeb.NewConfiguration()
	client := koyeb.NewAPIClient(cfg)

	ctx := context.Background()
	ctx = context.WithValue(ctx, koyeb.ContextAccessToken, c.Token)

	rqst := client.DomainsApi.ListDomains(ctx)
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
