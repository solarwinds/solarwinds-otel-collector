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

package e2e

import (
	"context"
	"log"
	"path/filepath"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	receivingContainer     = "receiver"
	testedContainer        = "sut"
	generatingContainer    = "generator"
	port                   = 17016
	collectorRunningPeriod = 35 * time.Second
	samplesCount           = 10
)

func runReceivingSolarWindsOTELCollector(
	ctx context.Context,
	networkName string,
) (testcontainers.Container, error) {
	configPath, err := filepath.Abs(filepath.Join(".", "testdata", "receiving_collector.yaml"))
	if err != nil {
		return nil, err
	}

	container, err := runSolarWindsOTELCollector(ctx, networkName, receivingContainer, configPath)
	return container, err
}

func runTestedSolarWindsOTELCollector(
	ctx context.Context,
	networkName string,
) (testcontainers.Container, error) {
	configPath, err := filepath.Abs(filepath.Join(".", "testdata", "emitting_collector.yaml"))
	if err != nil {
		return nil, err
	}

	container, err := runSolarWindsOTELCollector(ctx, networkName, testedContainer, configPath)
	return container, err
}

func runSolarWindsOTELCollector(
	ctx context.Context,
	networkName string,
	containerName string,
	configPath string,
) (testcontainers.Container, error) {
	lc := new(logConsumer)
	lc.Prefix = containerName
	req := testcontainers.ContainerRequest{
		Image: "solarwinds-otel-collector:latest",
		LogConsumerCfg: &testcontainers.LogConsumerConfig{
			Consumers: []testcontainers.LogConsumer{lc},
		},
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      configPath,
				ContainerFilePath: "/opt/default-config.yaml",
				FileMode:          0o440,
			},
		},
		WaitingFor: wait.ForLog("Everything is ready. Begin running and processing data."),
		Networks:   []string{networkName},
		Name:       containerName,
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	return container, err
}

func runGeneratorContainer(
	ctx context.Context,
	networkName string,
	cmd []string,
) (testcontainers.Container, error) {
	containerName := generatingContainer

	lc := new(logConsumer)
	lc.Prefix = containerName

	req := testcontainers.ContainerRequest{
		Image: "ghcr.io/open-telemetry/opentelemetry-collector-contrib/telemetrygen:latest",
		LogConsumerCfg: &testcontainers.LogConsumerConfig{
			Consumers: []testcontainers.LogConsumer{lc},
		},
		Networks: []string{networkName},
		Name:     containerName,
		Cmd:      cmd,
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	return container, err
}

type logConsumer struct {
	Prefix string
}

func (lc *logConsumer) Accept(l testcontainers.Log) {
	log.Printf("***%s: %s", lc.Prefix, string(l.Content))
}
