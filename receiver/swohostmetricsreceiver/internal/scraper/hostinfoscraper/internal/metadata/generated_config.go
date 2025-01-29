// Code generated by mdatagen. DO NOT EDIT.

package metadata

import (
	"go.opentelemetry.io/collector/confmap"
)

// MetricConfig provides common config for a particular metric.
type MetricConfig struct {
	Enabled bool `mapstructure:"enabled"`

	enabledSetByUser bool
}

func (ms *MetricConfig) Unmarshal(parser *confmap.Conf) error {
	if parser == nil {
		return nil
	}
	err := parser.Unmarshal(ms)
	if err != nil {
		return err
	}
	ms.enabledSetByUser = parser.IsSet("enabled")
	return nil
}

// MetricsConfig provides config for swohostmetricsreceiver/hostinfo metrics.
type MetricsConfig struct {
	SwoHostinfoFirewall       MetricConfig `mapstructure:"swo.hostinfo.firewall"`
	SwoHostinfoUptime         MetricConfig `mapstructure:"swo.hostinfo.uptime"`
	SwoHostinfoUserLastLogged MetricConfig `mapstructure:"swo.hostinfo.user.lastLogged"`
}

func DefaultMetricsConfig() MetricsConfig {
	return MetricsConfig{
		SwoHostinfoFirewall: MetricConfig{
			Enabled: false,
		},
		SwoHostinfoUptime: MetricConfig{
			Enabled: false,
		},
		SwoHostinfoUserLastLogged: MetricConfig{
			Enabled: false,
		},
	}
}

// MetricsBuilderConfig is a configuration for swohostmetricsreceiver/hostinfo metrics builder.
type MetricsBuilderConfig struct {
	Metrics MetricsConfig `mapstructure:"metrics"`
}

func DefaultMetricsBuilderConfig() MetricsBuilderConfig {
	return MetricsBuilderConfig{
		Metrics: DefaultMetricsConfig(),
	}
}
