//go:build integration

package e2e

import (
	"context"
	"github.com/testcontainers/testcontainers-go"
	"io"
	"log"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"

	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	receivingContainer  = "receiver"
	testedContainer     = "sut"
	generatingContainer = "generator"
	port                = 17016
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
		"--metrics", "10",
		"--otlp-insecure",
		"--otlp-endpoint", "sut:17016",
		"--otlp-attributes", "resource.attributes.testing_attribute=\"testing_value\"",
	}

	gContainer, err := runGeneratorContainer(ctx, net.Name, cmd)
	require.NoError(t, err)
	testcontainers.CleanupContainer(t, gContainer)

	<-time.After(10 * time.Second)
	log.Println("***: evaluation in progress")

	expectedMetricsCount := 10
	evaluateMetricsStream(t, ctx, rContainer, expectedMetricsCount)
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

	ms := pmetric.NewMetrics()
	jum := new(pmetric.JSONUnmarshaler)
	for _, line := range lines {
		m, err := jum.UnmarshalMetrics([]byte(line))
		if err != nil {
			continue
		}

		require.Equal(t, m.ResourceMetrics().Len(), 1, "it must contain exactly one resource metric")
		evaluateResourceAttributes(t, m.ResourceMetrics().At(0).Resource().Attributes())
		m.ResourceMetrics().MoveAndAppendTo(ms.ResourceMetrics())
	}
	require.Equal(t, ms.MetricCount(), expectedCount)
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
		"--traces", "10",
		"--otlp-insecure",
		"--otlp-endpoint", "sut:17016",
		"--otlp-attributes", "resource.attributes.testing_attribute=\"testing_value\"",
	}

	gContainer, err := runGeneratorContainer(ctx, net.Name, cmd)
	require.NoError(t, err)
	testcontainers.CleanupContainer(t, gContainer)

	<-time.After(10 * time.Second)
	log.Println("***: evaluation in progress")

	expectedTracesCount := 10 * 2
	evaluateTracesStream(t, ctx, rContainer, expectedTracesCount)
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
	jum := new(ptrace.JSONUnmarshaler)
	for _, line := range lines {
		tr, err := jum.UnmarshalTraces([]byte(line))
		if err != nil {
			continue
		}

		require.Equal(t, tr.ResourceSpans().Len(), 1, "it must contain exactly one resource span")
		evaluateResourceAttributes(t, tr.ResourceSpans().At(0).Resource().Attributes())
		tr.ResourceSpans().MoveAndAppendTo(trs.ResourceSpans())
	}
	require.Equal(t, expectedCount, trs.SpanCount())
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
		"--logs", "10",
		"--body", "testing log body",
		"--otlp-insecure",
		"--otlp-endpoint", "sut:17016",
		"--otlp-attributes", "resource.attributes.testing_attribute=\"testing_value\"",
	}

	gContainer, err := runGeneratorContainer(ctx, net.Name, cmd)
	require.NoError(t, err)
	testcontainers.CleanupContainer(t, gContainer)

	<-time.After(10 * time.Second)
	log.Println("***: evaluation in progress")

	expectedLogsCount := 10
	evaluateLogsStream(t, ctx, rContainer, expectedLogsCount)
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
	jum := new(plog.JSONUnmarshaler)
	for _, line := range lines {
		lg, err := jum.UnmarshalLogs([]byte(line))
		if err != nil {
			continue
		}

		require.Equal(t, lg.ResourceLogs().Len(), 1, "it must contain exactly one resource log")
		evaluateResourceAttributes(t, lg.ResourceLogs().At(0).Resource().Attributes())
		lg.ResourceLogs().MoveAndAppendTo(lgs.ResourceLogs())
	}
	require.Equal(t, expectedCount, lgs.LogRecordCount())
}

func evaluateResourceAttributes(
	t *testing.T,
	atts pcommon.Map,
) {
	val, ok := atts.Get("resource.attributes.testing_attribute")
	require.True(t, ok, "testing attribute must exist")
	require.Equal(t, val.AsString(), "testing_value", "testing attribute value must be the same")
}

func runReceivingSolarWindsOTELCollector(
	ctx context.Context,
	networkName string,
) (testcontainers.Container, error) {
	containerName := receivingContainer

	configPath, err := filepath.Abs(filepath.Join(".", "testdata", "receiving_collector.yaml"))
	if err != nil {
		return nil, err
	}

	lc := new(MyLogConsumer)
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

func runTestedSolarWindsOTELCollector(
	ctx context.Context,
	networkName string,
) (testcontainers.Container, error) {
	containerName := testedContainer

	configPath, err := filepath.Abs(filepath.Join(".", "testdata", "emitting_collector.yaml"))
	if err != nil {
		return nil, err
	}

	lc := new(MyLogConsumer)
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

	lc := new(MyLogConsumer)
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

	log.Print("*** raw content:\n" + string(content) + "\n")
	lines := strings.Split(string(content), "\n")
	return lines, nil
}

type MyLogConsumer struct {
	Prefix string
}

func (lc *MyLogConsumer) Accept(l testcontainers.Log) {
	log.Printf("***%s: %s", lc.Prefix, string(l.Content))
}
