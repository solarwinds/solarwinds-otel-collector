//go:build !integration

package lastloggeduser

import (
	"fmt"
	"testing"

	"github.com/solarwinds/solarwinds-otel-collector/receiver/swohostmetricsreceiver/internal/providers"
	"github.com/solarwinds/solarwinds-otel-collector/receiver/swohostmetricsreceiver/internal/providers/loggedusers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/solarwinds-cloud/uams-plugin-lib/pkg/logger"
)

func Test_Functional(t *testing.T) {
	t.Skip("This test should be run manually only")

	_ = logger.Setup(logger.WithLogToStdout(true))

	sut := NewEmitter()

	err := sut.Init()
	assert.Nil(t, err)

	er := sut.Emit()
	assert.Nil(t, er.Error)

	fmt.Printf("Result: %+v\n", er.Data)
}

func Test_Initialize_NotFailing(t *testing.T) {
	sut := NewEmitter()
	err := sut.Init()
	require.NoError(t, err)
}

func Test_GetEmittingFunction_emit_WhenReceivedDataEmitsMetric(t *testing.T) {
	usersProvider := &usersProviderMock{
		Data: loggedusers.Data{
			Users: []loggedusers.User{{
				Name:        `Test Name`,
				DisplayName: `Test DisplayName`,
			}},
		},
	}

	sut := createMetricEmitter(usersProvider)

	er := sut.Emit()
	assert.Nil(t, er.Error, "Emitter must not fail. Error:[%+v]", er.Error)

	metricCount := er.Data.Len()
	assert.Equal(t, 1, metricCount, "Expected number of metrics is 1")

	metric := er.Data.At(0)
	assert.Equal(t, MetricName, metric.Name())
	assert.Equal(t, MetricDescription, metric.Description())
	assert.Equal(t, MetricUnit, metric.Unit())

	points := metric.Gauge().DataPoints()
	assert.Equal(t, 1, points.Len(), "Metric count is different than expected")

	point := points.At(0)
	assert.Equal(t, int64(1), point.IntValue(), "Metric value is different than expected")

	attributes := point.Attributes()
	assert.Equal(t, 2, attributes.Len(), "Count of attributes is different than expected")
	expectedAttributes := map[string]any{
		"user.name":        "Test Name",
		"user.displayname": "Test DisplayName",
	}
	assert.EqualValues(t, expectedAttributes, attributes.AsRaw())
}

func Test_GetEmittingFunction_emit_WhenReceivedErrorEmitsEmptyMetricAndError(t *testing.T) {
	expectedError := fmt.Errorf("cardinal mistake")
	usersProvider := &usersProviderMock{
		Data: loggedusers.Data{Error: expectedError},
	}

	sut := createMetricEmitter(usersProvider)

	er := sut.Emit()
	assert.Equal(t, 0, er.Data.Len(), "Metric slice must be empty")
	assert.Equal(t, expectedError, er.Error, "Emitter must fail with %+v", expectedError)
}

type usersProviderMock struct {
	Data loggedusers.Data
}

var _ providers.Provider[loggedusers.Data] = (*usersProviderMock)(nil)

func (m *usersProviderMock) Provide() <-chan loggedusers.Data {
	ch := make(chan loggedusers.Data, 1)
	ch <- m.Data
	close(ch)
	return ch
}
