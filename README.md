# SolarWinds OpenTelemetry Collector
SolarWinds OpenTelemetry Collector is a distribution of OpenTelemetry Collector with components
bundled from [opentelemetry-collector](https://github.com/open-telemetry/opentelemetry-collector/tree/main)
and [opentelemetry-collector-contrib](https://github.com/open-telemetry/opentelemetry-collector-contrib). It also contains specific SolarWinds components for easier usage and enhanced telemetry collection.


## Installation
### Docker

Get the image from DockerHub.

`docker pull solarwinds/solarwinds-otel-collector:latest`

To run the image utilize following command:

`docker run  -v ./your_config_file.yaml:/opt/default-config.yaml solarwinds-otel-collector:latest`

See [complete docker documentation](./build/docker/README.md).

### Binary
Build the binary in `cmd/solarwinds-otel-collector` with `go build -tags full .`

After successful build, `solarwinds-otel-collector` should be present in `cmd/solarwinds-otel-collector`.

Run `solarwinds-otel-collector --config=config.yaml`.

## Getting Started
Configuration for SolarWinds OTel Collector has to contain [SolarWinds Extension](./extension/solarwindsextension/README.md) and [Solarwinds Exporter](./exporter/solarwindsexporter/README.md). 

### Example with docker
1. Generate your ingestion token in SWO. See [API Tokens](https://documentation.solarwinds.com/en/success_center/observability/content/settings/api-tokens.htm).
2. Create a `config.yaml` file that contains configuration for the SolarWinds OTel Collector. Insert the ingestion token and choose a correct data center (na-01, na-02, eu-01, ap-01). Specify the collector name.
```yaml
service:
  extensions: [solarwinds]
  pipelines:
    metrics:
      receivers: [redis]
      exporters: [solarwinds]
receivers:
  redis:
    endpoint: "<redis-url>:6379"
    collection_interval: 10s
    password: ${env:REDIS_PASSWORD}
extensions:
  solarwinds:
    token: "<ingestion-token>"
    data_center: "na-01"
    collector_name: "<collector-name>"

exporters:
  solarwinds:
```
3. Pull the SolarWinds OTel Collector from DockerHub.
```
docker pull solarwinds/solarwinds-otel-collector:latest
```
4. Start the container with your `config.yaml`. 
```
docker run  -v ./config.yaml:/opt/default-config.yaml solarwinds/solarwinds-otel-collector:latest
```

## Components
The SolarWinds OpenTelemetry collector contains following components:
- receivers
  - full set of [opentelemetry-collector receivers](https://github.com/open-telemetry/opentelemetry-collector/tree/v0.113.0/receiver)
  - full set of [opentelemetry-collector-contrib receivers](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/v0.113.0/receiver)
- processors
  - full set of [opentelemetry-collector processors](https://github.com/open-telemetry/opentelemetry-collector/tree/v0.113.0/processor)
  - full set of [opentelemetry-collector-contrib processors](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/v0.113.0/processor)
- exporters
    - [`solarwindsexporter`](./exporter/solarwindsexporter)
    - [`otlpexporter`](https://github.com/open-telemetry/opentelemetry-collector/tree/v0.113.0/exporter/otlpexporter)
    - [`fileexporter`](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/v0.113.0/exporter/fileexporter)
    - [`debugexporter`](https://github.com/open-telemetry/opentelemetry-collector/tree/v0.113.0/exporter/debugexporter)
    - [`nopexporter`](https://github.com/open-telemetry/opentelemetry-collector/tree/v0.113.0/exporter/nopexporter)
- extensions
    - [`solarwindsextension`](./extension/solarwindsextension)
    - full set of [opentelemetry-collector extensions](https://github.com/open-telemetry/opentelemetry-collector/tree/v0.113.0/extension)
    - full set of [opentelemetry-collector-contrib extensions](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/v0.113.0/extension)
- connectors
    - full set of [opentelemetry-collector connectors](https://github.com/open-telemetry/opentelemetry-collector/tree/v0.113.0/connector)
    - full set of [opentelemetry-collector-contrib connectors](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/v0.113.0/connector)

## Contributing
See [CONTRIBUTING.md](./CONTRIBUTING.md).
