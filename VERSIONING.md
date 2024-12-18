# Versioning
In general, versioning and stability standards of components maintained in this repository follow the standards laid out by [OpenTelemetry Collector](https://github.com/open-telemetry/opentelemetry-collector/blob/main/VERSIONING.md).

## SolarWinds OpenTelemetry Collector
The version of SolarWinds OpenTelemetry Collector and its distributions (i.e. Docker image) is based on the version of OpenTelemetry Collector consumed, with reservation of _patch_ version for feature set specific to this distribution.

As a result, the SolarWinds OpenTelemetry collector version is composed in following way:
* **Major** version follows major version of OpenTelemetry Collector consumed
* **Minor** version follows minor version of OpenTelemetry Collector consumed
* **Patch** version is reserved for SolarWinds-specific feature set.

_**Example**: SolarWinds OpenTelemetry Collector v0.113.2 consumes OpenTelemetry Collector 0.113.X and contains 3-distribution specific additions (i.e. updates or additions to components consumed)_
