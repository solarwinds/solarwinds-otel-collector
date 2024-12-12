// Copyright 2024 SolarWinds Worldwide, LLC. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Program solarwinds-otel-collector is an OpenTelemetry Collector binary.
package main

import (
	"log"

	"github.com/open-telemetry/opentelemetry-collector-contrib/confmap/provider/aesprovider"
	"github.com/open-telemetry/opentelemetry-collector-contrib/confmap/provider/s3provider"
	"github.com/open-telemetry/opentelemetry-collector-contrib/confmap/provider/secretsmanagerprovider"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/confmap/provider/envprovider"
	"go.opentelemetry.io/collector/confmap/provider/fileprovider"
	"go.opentelemetry.io/collector/confmap/provider/httpprovider"
	"go.opentelemetry.io/collector/confmap/provider/httpsprovider"
	"go.opentelemetry.io/collector/confmap/provider/yamlprovider"
	"go.opentelemetry.io/collector/otelcol"
)

func main() {
	info := component.BuildInfo{
		Command:     "solarwinds-otel-collector",
		Description: "SolarWinds distribution for OpenTelemetry",
		Version:     "0.113.0",
	}

	err := run(otelcol.CollectorSettings{
		BuildInfo: info,
		Factories: components,
		ConfigProviderSettings: otelcol.ConfigProviderSettings{
			ResolverSettings: confmap.ResolverSettings{
				ProviderFactories: []confmap.ProviderFactory{
					envprovider.NewFactory(),
					fileprovider.NewFactory(),
					httpprovider.NewFactory(),
					httpsprovider.NewFactory(),
					yamlprovider.NewFactory(),
					aesprovider.NewFactory(),
					s3provider.NewFactory(),
					secretsmanagerprovider.NewFactory(),
				},
			},
		},
	})

	if err != nil {
		log.Fatalf("collector server run finished with error: %v", err)
	}
}

func run(params otelcol.CollectorSettings) error {
	cmd := otelcol.NewCommand(params)
	err := cmd.Execute()
	return err
}
