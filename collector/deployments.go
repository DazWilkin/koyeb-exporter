package collector

import (
	"context"
	"log"

	"github.com/koyeb/koyeb-api-client-go/api/v1/koyeb"
	"github.com/prometheus/client_golang/prometheus"
)

// Ensure that DeploymentCollector implements Prometheus' Collector interface
var _ prometheus.Collector = (*DeploymentsCollector)(nil)

// DeploymentsCollector collects Koyeb Deployments metrics
type DeploymentsCollector struct {
	Ctx    context.Context
	Client *koyeb.APIClient

	Up *prometheus.Desc
}

// NewDeploymentsCollector is a function that creates a new DeploymentsCollector
func NewDeploymentsCollector(ctx context.Context, client *koyeb.APIClient) *DeploymentsCollector {
	subsystem := "deployments"
	return &DeploymentsCollector{
		Ctx:    ctx,
		Client: client,
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
	rqst := c.Client.DeploymentsApi.ListDeployments(c.Ctx)
	resp, _, err := rqst.Execute()
	if err != nil {
		msg := "unable to list Deployments"
		log.Printf(msg, err)
		return
	}

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
