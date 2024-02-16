package honeycombeventsreceiver

import (
	"context"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/config/confignet"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/internal/localhostgate"
	"go.opentelemetry.io/collector/receiver"
)

const (
	httpPort = 4318
)

// NewFactory creates a factory for the StatsD receiver.
func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		metadata.Type,
		createDefaultConfig,
		receiver.WithTraces(createTracesReceiver, component.StabilityLevelDevelopment),
		receiver.WithLogs(createLogsReceiver, component.StabilityLevelDevelopment),
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		Protocols: Protocols{
			HTTP: &HTTPConfig{
				ServerConfig: &confighttp.ServerConfig{
					Endpoint: localhostgate.EndpointForPort(httpPort),
				},
			},
		},
	}
}

func createTracesReceiver(
	_ context.Context,
	set receiver.CreateSettings,
	cfg component.Config,
	nextConsumer consumer.Traces,
) (receiver.Traces, error) {
	var err error
	r := receivers.GetOrAdd(cfg, func() component.Component {
		rCfg := cfg.(*Config)
		var recv *ocReceiver
		recv, err = newOpenCensusReceiver(rCfg.NetAddr.Transport, rCfg.NetAddr.Endpoint, nil, nil, set, rCfg.buildOptions()...)
		return recv
	})
	if err != nil {
		return nil, err
	}
	r.Unwrap().(*ocReceiver).traceConsumer = nextConsumer

	return r, nil
}

func createLogsReceiver(
	_ context.Context,
	set receiver.CreateSettings,
	cfg component.Config,
	nextConsumer consumer.Logs,
) (receiver.Logs, error) {
	var err error
	r := receivers.GetOrAdd(cfg, func() component.Component {
		rCfg := cfg.(*Config)
		var recv *ocReceiver
		recv, err = newOpenCensusReceiver(rCfg.NetAddr.Transport, rCfg.NetAddr.Endpoint, nil, nil, set, rCfg.buildOptions()...)
		return recv
	})
	if err != nil {
		return nil, err
	}
	r.Unwrap().(*ocReceiver).metricsConsumer = nextConsumer

	return r, nil
}
