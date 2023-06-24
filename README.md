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
|`apps_up`|Gauge|1 if the App is up, 0 otherwise|
|`credentials_up`|Gauge|1 if the Credential is up, 0 otherwise|
|`deployments_up`|Gauge|1 if the Deployment is up, 0 otherwise|
|`domains_up`|Gauge|1 if the Domain is up, 0 otherwise|
|`exporter_build_info`|Counter|A metric with a constant '1' value labeled by OS version, Go version, and the Git commit of the exporter|
|`exporter_start_time`|Gauge|Exporter start time in Unix epoch seconds|
|`instances_up`|Gauge|1 if the instance is up, 0 otherwise|
|`services_up`|Gauge|1 if the Service is up, 0 otherwise|

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
  - alert: koyeb_credentials_up
    expr: min_over_time(koyeb_credentials_up{}[15m]) > 0
    for: 3h
    labels:
      severity: page
    annotations:
      summary: "Koyeb Credentials ({{ $value }}) up (name: {{ $labels.name }})"
  - alert: koyeb_deployments_up
    expr: min_over_time(koyeb_deployments_up{}[15m]) > 0
    for: 3h
    labels:
      severity: page
    annotations:
      summary: "Koyeb Deployments ({{ $value }}) up (name: {{ $labels.name }})"
  - alert: koyeb_domains_up
    expr: min_over_time(koyeb_domains_up{}[15m]) > 0
    for: 3h
    labels:
      severity: page
    annotations:
      summary: "Koyeb Domains ({{ $value }}) up (name: {{ $labels.name }})"
  - alert: koyeb_instances_up
    expr: min_over_time(koyeb_instances_up{}[15m]) > 0
    for: 3h
    labels:
      severity: page
    annotations:
      summary: "Koyeb Instances ({{ $value }}) up (region: {{ $labels.region }})"
  - alert: koyeb_secrets_up
    expr: min_over_time(koyeb_secrets_up{}[15m]) > 0
    for: 3h
    labels:
      severity: page
    annotations:
      summary: "Koyeb Secrets ({{ $value }}) up (name: {{ $labels.name }})"
```

<hr/>
<br/>
<a href="https://www.buymeacoffee.com/dazwilkin" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/default-orange.png" alt="Buy Me A Coffee" height="41" width="174"></a>