# E2E Tests
The purpose of the end-to-end tests is to validate the overall integrity of the SolarWinds OpenTelemetry Collector distribution and the components within it, ensuring the correct handling of all supported signals throughout the pipeline — from receiving signals to their publishing.

## How to Develop
For local test runs, it is required to have a local build of the SolarWinds OpenTelemetry Collector tagged as `solarwinds-otel-collector:latest`.

To run the tests, execute the `make e2e-tests` command.

> [!NOTE]  
> Please refer to the documentation on how to produce the local Docker image [here](../../build/docker/README.md).

## Technology Stack
The tests utilize the following technologies:

* [testcontainers](https://testcontainers.com/?language=go) - to create running containers utilizing the SolarWinds OpenTelemetry Collector at test runtime
* [telemetrygen tool](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/cmd/telemetrygen) - to generate telemetry data sent to the receiving OTLP endpoint
    * **NOTE:** telemetrygen is consumed as a [Docker image](https://github.com/open-telemetry/opentelemetry-collector-contrib/pkgs/container/opentelemetry-collector-contrib%2Ftelemetrygen) and run as a separate container during test runtime
