# Prometheus Exporter for [Koyeb](https://koyeb.com)

[![build-container](https://github.com/DazWilkin/koyeb-exporter/actions/workflows/build.yml/badge.svg)](https://github.com/DazWilkin/koyeb-exporter/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/DazWilkin/koyeb-exporter)](https://goreportcard.com/report/github.com/DazWilkin/koyeb-exporter)

+ `ghcr.io/dazwilkin/koyeb-exporter:1252430839889e215ae64f9819afedb07aa1ce0a`

Exports Koyeb (Apps, Deployments, Instances) to enable e.g. (Prometheus) Alerting on Koyeb resource consumption ($$$).

## Run

[`koyeb-exporter`](https://github.com/DazWilkin/koyeb-exporter/pkgs/container/koyeb-exporter)

```bash
TOKEN=$(more ~/.koyeb.yaml | yq .token) # Koyeb API Token

PORT="..."

podman run \
--interactive --tty --rm \
--env=TOKEN=${TOKEN} \
ghcr.io/dazwilkin/koyeb-exporter:1252430839889e215ae64f9819afedb07aa1ce0a \
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
|`secrets_up`|Gauge|1 if the Secret is up, 0 otherwise|
|`services_up`|Gauge|1 if the Service is up, 0 otherwise|

## Prometheus

```bash
VERS="v2.45.0"

# Binds to host network to scrape Koyeb Exporter
podman run \
--interactive --tty --rm \
--net=host \
--volume=${PWD}/prometheus.yml:/etc/prometheus/prometheus.yml \
--volume=${PWD}/rules.yml:/etc/alertmanager/rules.yml \
quay.io/prometheus/prometheus:${VERS} \
  --config.file=/etc/prometheus/prometheus.yml \
  --web.enable-lifecycle
```

See [`prometheus.yml`](/prometheus.yml)

## Alerting Rules

See [`rules.yml`](/rules.yml)

## Sigstore

`koyeb-exporter`` container images are being signed by Sigstore and may be verified:
```bash
cosign verify \
--key=./cosign.pub \
ghcr.io/dazwilkin/koyeb-exporter:1252430839889e215ae64f9819afedb07aa1ce0a
```

> **NOTE** `cosign.pub` may be downloaded [here](./cosign.pub)

To install cosign:
```bash
go install github.com/sigstore/cosign/cmd/cosign@latest
```

<hr/>
<br/>
<a href="https://www.buymeacoffee.com/dazwilkin" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/default-orange.png" alt="Buy Me A Coffee" height="41" width="174"></a>