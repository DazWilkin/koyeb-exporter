package collector

import (
	"context"
	"log"

	"github.com/DazWilkin/koyeb-exporter/types"
	"github.com/koyeb/koyeb-api-client-go/api/v1/koyeb"
	"github.com/prometheus/client_golang/prometheus"
)

// Ensure that DeploymentCollector implements Prometheus' Collector interface
var _ prometheus.Collector = (*SecretsCollector)(nil)

// SecretsCollector collects Koyeb Secrets metrics
type SecretsCollector struct {
	Token string

	Up *prometheus.Desc
}

// NewSecretsCollector is a function that creates a new SecretsCollector
func NewSecretsCollector(token string) *SecretsCollector {
	subsystem := "secrets"
	return &SecretsCollector{
		Token: token,

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
	cfg := koyeb.NewConfiguration()
	client := koyeb.NewAPIClient(cfg)

	ctx := context.Background()
	ctx = context.WithValue(ctx, koyeb.ContextAccessToken, c.Token)

	rqst := client.SecretsApi.ListSecrets(ctx)
	resp, _, err := rqst.Execute()
	if err != nil {
		msg := "unable to list Secrets"
		log.Printf(msg, err)
		return
	}

	for _, secret := range resp.Secrets {
		log.Println(secret)
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
