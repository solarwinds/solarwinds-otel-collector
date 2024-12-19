# SolarWinds OpenTelemetry Collector
Distribution of OpenTelemetry Collector with all available components
bundled within from [opentelemetry-collector](https://github.com/open-telemetry/opentelemetry-collector/tree/main)
and [opentelemetry-collector-contrib](https://github.com/open-telemetry/opentelemetry-collector-contrib).

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

## Installation
### Binary
- build in `cmd/solarwinds-otel-collector` with `go build .`
- after successful build, `solarwinds-otel-collector` should be present in `cmd/solarwinds-otel-collector`
- run `solarwinds-otel-collector --config=example_config.yaml`

### Docker
See full [docker documentation](./build/docker/README.md).

### CI/CD
The _SolarWinds OpenTelemetry Collector_ utilizes [GitHub Actions pipeline](./.github). 
The standard build pipeline is triggered with each PR opened to main or release branch and consists of:
- build of _SolarWinds OpenTelemetry Collector_ docker image

The release pipeline is triggered with designated tag publishing and consists of:
- build of _SolarWinds OpenTelemetry Collector_ docker image and its publishing to DockerHub
- creation of GitHub release

## Contributing
See [CONTRIBUTING.md](./CONTRIBUTING.md).
