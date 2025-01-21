package scraper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"

	"github.com/solarwinds/solarwinds-otel-collector/receiver/swohostmetricsreceiver/internal/types"
)

// RequireComponentType converts a `from` string
// to `component.Type` and fails the test if it cannot
// be converted.
func RequireComponentType(t *testing.T, from string) component.Type {
	ty, err := component.NewType(from)
	require.NoError(t, err)

	return ty
}

type DummyMetricEmitterMock struct{}

var _ types.MetricEmitterInterface = (*DummyMetricEmitterMock)(nil)

// GetEmittingFunction implements types.MetricEmitterInterface.
func (*DummyMetricEmitterMock) GetEmittingFunction() types.MetricsEmittingFunc {
	return func() (pmetric.MetricSlice, error) {
		return pmetric.NewMetricSlice(), nil
	}
}

// Initialize implements types.MetricEmitterInterface.
func (*DummyMetricEmitterMock) Initialize() error {
	return nil
}

type metricEmitterMock struct {
	metricSlice pmetric.MetricSlice
}

var _ types.MetricEmitterInterface = (*metricEmitterMock)(nil)

func CreateMetricEmitterMock(ms pmetric.MetricSlice) types.MetricEmitterInterface {
	return &metricEmitterMock{
		metricSlice: ms,
	}
}

// GetEmittingFunction implements types.MetricEmitterInterface.
func (m *metricEmitterMock) GetEmittingFunction() types.MetricsEmittingFunc {
	return func() (pmetric.MetricSlice, error) {
		return m.metricSlice, nil
	}
}

// Initialize implements types.MetricEmitterInterface.
func (*metricEmitterMock) Initialize() error {
	return nil
}

func CreateArtificialMetricSlice(
	metricName string,
	numericDpValues []int64,
) pmetric.MetricSlice {
	ms := pmetric.NewMetricSlice()

	m := ms.AppendEmpty()
	m.SetDescription("Artificial metric providing constant value in seconds")
	m.SetName(metricName)
	m.SetUnit("s")

	dps := m.SetEmptySum().DataPoints()
	for _, v := range numericDpValues {
		dp := dps.AppendEmpty()
		dp.SetIntValue(v)
		dp.SetStartTimestamp(pcommon.NewTimestampFromTime(time.Now()))
		dp.SetTimestamp(pcommon.NewTimestampFromTime(time.Now()))
	}

	return ms
}

type failingEmitterMock struct {
	InitError error
	EmitError error
}

var _ types.MetricEmitterInterface = (*failingEmitterMock)(nil)

func CreateFailingEmitterMock(
	ie error,
	ee error,
) types.MetricEmitterInterface {
	return &failingEmitterMock{
		InitError: ie,
		EmitError: ee,
	}
}

// Initialize implements types.MetricEmitterInterface.
func (e *failingEmitterMock) Initialize() error {
	return e.InitError
}

// GetEmittingFunction implements types.MetricEmitterInterface.
func (e *failingEmitterMock) GetEmittingFunction() types.MetricsEmittingFunc {
	return func() (pmetric.MetricSlice, error) {
		return pmetric.NewMetricSlice(), e.EmitError
	}
}
