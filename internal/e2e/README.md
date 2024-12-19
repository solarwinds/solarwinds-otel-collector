# E2E Tests
The purpose of the end-to-end tests is to validate the overall integrity of SolarWinds OpenTelemetry Collector distribution and the components used within, verifying correct handling of all supported signals throughout the pipeline - from signals receiving to their publishing.

## How to Develop
For local test runs it is required to have local build of SolarWinds OpenTelemetry Collector tagged as `solarwinds-otel-collector:latest`.

To run the tests, call `make e2e-tests` command

> [!NOTE]  
> Please refer to documentation on how to produce the local docker image [here](../../build/docker/README.md).
## Technology Stack
Tests utilize following technologies within:
* [testcontainers](https://testcontainers.com/?language=go) - to create running containers utilizing SolarWinds OpenTelemetry Collector at test runtime
* [telemetrygen tool](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/cmd/telemetrygen) - to generate telemetry data sent to receiving OTLP endpoint
    * **NOTE:** telemetrygen is consumed in form of [docker image](https://github.com/open-telemetry/opentelemetry-collector-contrib/pkgs/container/opentelemetry-collector-contrib%2Ftelemetrygen) and run as separate container at tests runtime
