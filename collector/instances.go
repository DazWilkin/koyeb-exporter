package collector

import (
	"context"
	"log"

	"github.com/koyeb/koyeb-api-client-go/api/v1/koyeb"
	"github.com/prometheus/client_golang/prometheus"
)

// Ensure that InstancesCollector implements Prometheus' Collector interface
var _ prometheus.Collector = (*InstancesCollector)(nil)

// InstancesCollector collects Koyeb Apps metrics
type InstancesCollector struct {
	Token string

	Up *prometheus.Desc
}

// NewInstancesCollector is a function that creates a new InstancesCollector
func NewInstancesCollector(token string) *InstancesCollector {
	subsystem := "instances"
	return &InstancesCollector{
		Token: token,

		Up: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "up"),
			"1 if the instance is up, 0 otherwise",
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
	cfg := koyeb.NewConfiguration()
	client := koyeb.NewAPIClient(cfg)

	ctx := context.Background()
	ctx = context.WithValue(ctx, koyeb.ContextAccessToken, c.Token)

	rqst := client.InstancesApi.ListInstances(ctx)
	resp, _, err := rqst.Execute()
	if err != nil {
		msg := "unable to list instances"
		log.Printf(msg, err)
		return
	}

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
