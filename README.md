# Prometheus Exporter for [Koyeb](https://koyeb.com)

Exports Koyeb (Apps, Deployments, Instances) to enable e.g. (Prometheus) Alerting on Koyeb resource consumption ($$$).

## Run

[`koyeb-exporter`](https://github.com/DazWilkin/koyeb-exporter/pkgs/container/koyeb-exporter)

```bash
TOKEN=$(more ~/.koyeb.yaml | yq .token) # Koyeb API Token
VERS="..." # See package list above
PORT="..."

podman run \
--interactive --tty --rm \
--env=TOKEN=${TOKEN} \
ghcr.io/dazwilkin/koyeb-exporter:${VERS} \
--endpoint=":${PORT} \
--path=/metrics
```

## Metrics

All metric names are prefix `koyeb_`

|Name|Type|Description|
|----|----|-----------|
|`apps_up`|Gauge|1 if the app is up, 0 otherwise|
|`deployments_up`|Gauge|1 if the deployment is up, 0 otherwise|
|`exporter_build_info`|Counter|A metric with a constant '1' value labeled by OS version, Go version, and the Git commit of the exporter|
|`exporter_start_time`|Gauge|Exporter start time in Unix epoch seconds|
|`instances_up`|Gauge|1 if the instance is up, 0 otherwise|

## Alerting Rules

```YAML
groups:
- name: koyeb_exporter
  rules:
  - alert: koyeb_apps_up
    expr: min_over_time(koyeb_apps_up{}[15m]) > 0
    for: 3h
    labels:
      severity: page
    annotations:
      summary: "Koyeb Apps ({{ $value }}) up (name: {{ $labels.name }})"
  - alert: koyeb_instances_up
    expr: min_over_time(koyeb_instances_up{}[15m]) > 0
    for: 3h
    labels:
      severity: page
    annotations:
      summary: "Koyeb Instances ({{ $value }}) up (region: {{ $labels.region }})"
```