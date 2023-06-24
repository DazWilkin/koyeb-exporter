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
	Token string

	Up *prometheus.Desc
}

// NewCredentialsCollector is a function that creates a new CredentialsCollector
func NewCredentialsCollector(token string) *CredentialsCollector {
	subsystem := "credentials"
	return &CredentialsCollector{
		Token: token,

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
	cfg := koyeb.NewConfiguration()
	client := koyeb.NewAPIClient(cfg)

	ctx := context.Background()
	ctx = context.WithValue(ctx, koyeb.ContextAccessToken, c.Token)

	rqst := client.CredentialsApi.ListCredentials(ctx)
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
