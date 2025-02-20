// Code generated by mdatagen. DO NOT EDIT.

package metadata

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/scraper/scrapertest"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

type testDataSet int

const (
	testDataSetDefault testDataSet = iota
	testDataSetAll
	testDataSetNone
)

func TestMetricsBuilder(t *testing.T) {
	tests := []struct {
		name        string
		metricsSet  testDataSet
		resAttrsSet testDataSet
		expectEmpty bool
	}{
		{
			name: "default",
		},
		{
			name:        "all_set",
			metricsSet:  testDataSetAll,
			resAttrsSet: testDataSetAll,
		},
		{
			name:        "none_set",
			metricsSet:  testDataSetNone,
			resAttrsSet: testDataSetNone,
			expectEmpty: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := pcommon.Timestamp(1_000_000_000)
			ts := pcommon.Timestamp(1_000_001_000)
			observedZapCore, observedLogs := observer.New(zap.WarnLevel)
			settings := scrapertest.NewNopSettings()
			settings.Logger = zap.New(observedZapCore)
			mb := NewMetricsBuilder(loadMetricsBuilderConfig(t, tt.name), settings, WithStartTime(start))

			expectedWarnings := 0

			assert.Equal(t, expectedWarnings, observedLogs.Len())

			defaultMetricsCount := 0
			allMetricsCount := 0

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordSwoHardwareinventoryCPUDataPoint(ts, 1, "processor.name-val", "processor.caption-val", "processor.manufacturer-val", "processor.model-val", "processor.stepping-val", "processor.cores-val", "processor.threads-val")

			res := pcommon.NewResource()
			metrics := mb.Emit(WithResource(res))

			if tt.expectEmpty {
				assert.Equal(t, 0, metrics.ResourceMetrics().Len())
				return
			}

			assert.Equal(t, 1, metrics.ResourceMetrics().Len())
			rm := metrics.ResourceMetrics().At(0)
			assert.Equal(t, res, rm.Resource())
			assert.Equal(t, 1, rm.ScopeMetrics().Len())
			ms := rm.ScopeMetrics().At(0).Metrics()
			if tt.metricsSet == testDataSetDefault {
				assert.Equal(t, defaultMetricsCount, ms.Len())
			}
			if tt.metricsSet == testDataSetAll {
				assert.Equal(t, allMetricsCount, ms.Len())
			}
			validatedMetrics := make(map[string]bool)
			for i := 0; i < ms.Len(); i++ {
				switch ms.At(i).Name() {
				case "swo.hardwareinventory.cpu":
					assert.False(t, validatedMetrics["swo.hardwareinventory.cpu"], "Found a duplicate in the metrics slice: swo.hardwareinventory.cpu")
					validatedMetrics["swo.hardwareinventory.cpu"] = true
					assert.Equal(t, pmetric.MetricTypeGauge, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Gauge().DataPoints().Len())
					assert.Equal(t, "CPU current clock speed in MHz.", ms.At(i).Description())
					assert.Equal(t, "MHz", ms.At(i).Unit())
					dp := ms.At(i).Gauge().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
					attrVal, ok := dp.Attributes().Get("processor.name")
					assert.True(t, ok)
					assert.EqualValues(t, "processor.name-val", attrVal.Str())
					attrVal, ok = dp.Attributes().Get("processor.caption")
					assert.True(t, ok)
					assert.EqualValues(t, "processor.caption-val", attrVal.Str())
					attrVal, ok = dp.Attributes().Get("processor.manufacturer")
					assert.True(t, ok)
					assert.EqualValues(t, "processor.manufacturer-val", attrVal.Str())
					attrVal, ok = dp.Attributes().Get("processor.model")
					assert.True(t, ok)
					assert.EqualValues(t, "processor.model-val", attrVal.Str())
					attrVal, ok = dp.Attributes().Get("processor.stepping")
					assert.True(t, ok)
					assert.EqualValues(t, "processor.stepping-val", attrVal.Str())
					attrVal, ok = dp.Attributes().Get("processor.cores")
					assert.True(t, ok)
					assert.EqualValues(t, "processor.cores-val", attrVal.Str())
					attrVal, ok = dp.Attributes().Get("processor.threads")
					assert.True(t, ok)
					assert.EqualValues(t, "processor.threads-val", attrVal.Str())
				}
			}
		})
	}
}
