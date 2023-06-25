package collector

import (
	"context"
	"log"

	"github.com/koyeb/koyeb-api-client-go/api/v1/koyeb"
	"github.com/prometheus/client_golang/prometheus"
)

// Ensure that DeploymentCollector implements Prometheus' Collector interface
var _ prometheus.Collector = (*CredentialsCollector)(nil)

// CredentialsCollector collects Koyeb Credentials metrics
type CredentialsCollector struct {
	Ctx    context.Context
	Client *koyeb.APIClient

	Up *prometheus.Desc
}

// NewCredentialsCollector is a function that creates a new CredentialsCollector
func NewCredentialsCollector(ctx context.Context, client *koyeb.APIClient) *CredentialsCollector {
	subsystem := "credentials"
	return &CredentialsCollector{
		Ctx:    ctx,
		Client: client,

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
	rqst := c.Client.CredentialsApi.ListCredentials(c.Ctx)
	resp, _, err := rqst.Execute()
	if err != nil {
		msg := "unable to list Credentials"
		log.Printf(msg, err)
		return
	}

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
