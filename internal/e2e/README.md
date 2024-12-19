# Integration Tests
The purpose of the integration tests is to validate the overall integrity of SolarWinds OpenTelemetry Collector distribution and the components used within, verifying correct handling of all supported signals.

## How to Develop
For local test runs it is required to have local build of SolarWinds OpenTelemetry Collector tagged as `solarwinds-otel-collector:latest`.

Tests then can be run with `go test -tags=integration` command.

> [NOTE!]  
> Please refer to documentation on how to produce the local docker image [here](../../build/docker/README.md).

## Technology Stack
Integration tests utilize following technologies within:
* [testcontainers](https://testcontainers.com/?language=go) - to create running containers utilizing SolarWinds OpenTelemetry Collector at test runtime
* [telemetrygen tool](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/cmd/telemetrygen) - to generate telemetry data sent to receiving OTLP endpoint
  * **NOTE:** telemetrygen is consumed in form of [docker image](https://github.com/open-telemetry/opentelemetry-collector-contrib/pkgs/container/opentelemetry-collector-contrib%2Ftelemetrygen) and run as separate container at tests runtime
