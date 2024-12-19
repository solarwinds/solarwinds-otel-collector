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

//go:build e2e

package e2e

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
)

const (
	resourceAttributeName  = "resource.attributes.testing_attribute"
	resourceAttributeValue = "testing_value"
)

func TestMetricStream(t *testing.T) {
	ctx := context.Background()

	net, err := network.New(ctx)
	require.NoError(t, err)
	testcontainers.CleanupNetwork(t, net)

	rContainer, err := runReceivingSolarWindsOTELCollector(ctx, net.Name)
	require.NoError(t, err)
	testcontainers.CleanupContainer(t, rContainer)

	eContainer, err := runTestedSolarWindsOTELCollector(ctx, net.Name)
	require.NoError(t, err)
	testcontainers.CleanupContainer(t, eContainer)

	cmd := []string{
		"metrics",
		"--metrics", strconv.Itoa(samplesCount),
		"--otlp-insecure",
		"--otlp-endpoint", fmt.Sprintf("%s:%d", testedContainer, port),
		"--otlp-attributes", fmt.Sprintf("%s=\"%s\"", resourceAttributeName, resourceAttributeValue),
	}

	gContainer, err := runGeneratorContainer(ctx, net.Name, cmd)
	require.NoError(t, err)
	testcontainers.CleanupContainer(t, gContainer)

	<-time.After(collectorRunningPeriod)

	evaluateMetricsStream(t, ctx, rContainer, samplesCount)
}

func TestTracesStream(t *testing.T) {
	ctx := context.Background()

	net, err := network.New(ctx)
	require.NoError(t, err)
	testcontainers.CleanupNetwork(t, net)

	rContainer, err := runReceivingSolarWindsOTELCollector(ctx, net.Name)
	require.NoError(t, err)
	testcontainers.CleanupContainer(t, rContainer)

	eContainer, err := runTestedSolarWindsOTELCollector(ctx, net.Name)
	require.NoError(t, err)
	testcontainers.CleanupContainer(t, eContainer)

	cmd := []string{
		"traces",
		"--traces", strconv.Itoa(samplesCount),
		"--otlp-insecure",
		"--otlp-endpoint", fmt.Sprintf("%s:%d", testedContainer, port),
		"--otlp-attributes", fmt.Sprintf("%s=\"%s\"", resourceAttributeName, resourceAttributeValue),
	}

	gContainer, err := runGeneratorContainer(ctx, net.Name, cmd)
	require.NoError(t, err)
	testcontainers.CleanupContainer(t, gContainer)

	<-time.After(collectorRunningPeriod)

	// Traces coming in couples.
	expectedTracesCount := samplesCount * 2
	evaluateTracesStream(t, ctx, rContainer, expectedTracesCount)
}

func TestLogsStream(t *testing.T) {
	ctx := context.Background()

	net, err := network.New(ctx)
	require.NoError(t, err)
	testcontainers.CleanupNetwork(t, net)

	rContainer, err := runReceivingSolarWindsOTELCollector(ctx, net.Name)
	require.NoError(t, err)
	testcontainers.CleanupContainer(t, rContainer)

	eContainer, err := runTestedSolarWindsOTELCollector(ctx, net.Name)
	require.NoError(t, err)
	testcontainers.CleanupContainer(t, eContainer)

	cmd := []string{
		"logs",
		"--logs", strconv.Itoa(samplesCount),
		"--otlp-insecure",
		"--otlp-endpoint", fmt.Sprintf("%s:%d", testedContainer, port),
		"--otlp-attributes", fmt.Sprintf("%s=\"%s\"", resourceAttributeName, resourceAttributeValue),
	}

	gContainer, err := runGeneratorContainer(ctx, net.Name, cmd)
	require.NoError(t, err)
	testcontainers.CleanupContainer(t, gContainer)

	<-time.After(collectorRunningPeriod)

	evaluateLogsStream(t, ctx, rContainer, samplesCount)
}

func evaluateMetricsStream(
	t *testing.T,
	ctx context.Context,
	container testcontainers.Container,
	expectedCount int,
) {
	// Obtain result from container.
	lines, err := loadResultFile(ctx, container, "/tmp/result.json")
	require.NoError(t, err)

	gms := pmetric.NewMetrics()
	hbms := pmetric.NewMetrics()
	jum := new(pmetric.JSONUnmarshaler)
	for _, line := range lines {
		m, err := jum.UnmarshalMetrics([]byte(line))
		if err != nil || m.ResourceMetrics().Len() == 0 {
			continue
		}

		if m.ResourceMetrics().At(0).ScopeMetrics().Len() == 0 ||
			m.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().Len() == 0 {
			continue
		}

		heartbeatMetricName := "sw.otelcol.uptime"
		generatedMetricName := "gen"
		metricName := m.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Name()

		switch metricName {
		case generatedMetricName:
			evaluateResourceAttributes(t, m.ResourceMetrics().At(0).Resource().Attributes())
			m.ResourceMetrics().MoveAndAppendTo(gms.ResourceMetrics())
		case heartbeatMetricName:
			m.ResourceMetrics().MoveAndAppendTo(hbms.ResourceMetrics())
		default:
			continue
		}
	}
	require.Equal(t, gms.MetricCount(), expectedCount)
	evaluateHeartbeatMetric(t, hbms)
}

func evaluateTracesStream(
	t *testing.T,
	ctx context.Context,
	container testcontainers.Container,
	expectedCount int,
) {
	// Obtain result from container.
	lines, err := loadResultFile(ctx, container, "/tmp/result.json")
	require.NoError(t, err)

	trs := ptrace.NewTraces()
	ms := pmetric.NewMetrics()
	tum := new(ptrace.JSONUnmarshaler)
	mum := new(pmetric.JSONUnmarshaler)
	for _, line := range lines {
		// Traces to process.
		tr, err := tum.UnmarshalTraces([]byte(line))
		if err == nil && tr.ResourceSpans().Len() != 0 {
			evaluateResourceAttributes(t, tr.ResourceSpans().At(0).Resource().Attributes())
			tr.ResourceSpans().MoveAndAppendTo(trs.ResourceSpans())
			continue
		}

		// Metrics to process.
		m, err := mum.UnmarshalMetrics([]byte(line))
		if err == nil && m.ResourceMetrics().Len() != 0 {
			m.ResourceMetrics().MoveAndAppendTo(ms.ResourceMetrics())
			continue
		}
	}

	evaluateHeartbeatMetric(t, ms)
	require.Equal(t, expectedCount, trs.SpanCount())
}

func evaluateLogsStream(
	t *testing.T,
	ctx context.Context,
	container testcontainers.Container,
	expectedCount int,
) {
	// Obtain result from container.
	lines, err := loadResultFile(ctx, container, "/tmp/result.json")
	require.NoError(t, err)

	lgs := plog.NewLogs()
	ms := pmetric.NewMetrics()
	lum := new(plog.JSONUnmarshaler)
	mum := new(pmetric.JSONUnmarshaler)
	for _, line := range lines {
		// Logs to process.
		lg, err := lum.UnmarshalLogs([]byte(line))
		if err == nil && lg.ResourceLogs().Len() != 0 {
			evaluateResourceAttributes(t, lg.ResourceLogs().At(0).Resource().Attributes())
			lg.ResourceLogs().MoveAndAppendTo(lgs.ResourceLogs())
			continue
		}

		// Metrics to process.
		m, err := mum.UnmarshalMetrics([]byte(line))
		if err == nil && m.ResourceMetrics().Len() != 0 {
			m.ResourceMetrics().MoveAndAppendTo(ms.ResourceMetrics())
			continue
		}
	}

	evaluateHeartbeatMetric(t, ms)
	require.Equal(t, expectedCount, lgs.LogRecordCount())
}

func evaluateHeartbeatMetric(
	t *testing.T,
	ms pmetric.Metrics,
) {
	require.GreaterOrEqual(t, ms.ResourceMetrics().Len(), 1, "there must be at least one metric")
	atts := ms.ResourceMetrics().At(0).Resource().Attributes()
	v, available := atts.Get("sw.otelcol.collector.name")
	require.True(t, available, "sw.otelcol.collector.name resource attribute must be available")
	require.Equal(t, "testing_collector_name", v.AsString(), "attribute value must be the same")
}

func evaluateResourceAttributes(
	t *testing.T,
	atts pcommon.Map,
) {
	val, ok := atts.Get(resourceAttributeName)
	require.True(t, ok, "testing attribute must exist")
	require.Equal(t, val.AsString(), resourceAttributeValue, "testing attribute value must be the same")
}

func loadResultFile(
	ctx context.Context,
	container testcontainers.Container,
	resultFilePath string,
) ([]string, error) {
	r, err := container.CopyFileFromContainer(ctx, resultFilePath)
	if err != nil {
		return make([]string, 0), err
	}

	content, err := io.ReadAll(r)
	if err != nil {
		return make([]string, 0), err
	}

	log.Print("*** raw result content:\n" + string(content) + "\n")
	lines := strings.Split(string(content), "\n")
	return lines, nil
}
