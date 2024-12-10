package internal

import (
	"context"
	"errors"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
	"sync"
	"time"
)

type metricsPusher func(ctx context.Context, md pmetric.Metrics) error
type metricsAdder func(ctx context.Context, md pmetric.Metrics) error

type Heartbeat struct {
	logger *zap.Logger

	cancel           context.CancelFunc
	startShutdownMtx sync.Mutex

	pushMetrics metricsPusher
	addMetrics  metricsAdder
}

var alreadyRunningError = errors.New("heartbeat already running")
var notRunningError = errors.New("heartbeat not started")

func NewHeartbeat(logger *zap.Logger, pushMetrics metricsPusher, addMetrics metricsAdder) *Heartbeat {
	logger.Debug("Creating Heartbeat")
	return &Heartbeat{logger: logger, pushMetrics: pushMetrics, addMetrics: addMetrics}
}

func (h *Heartbeat) Start() error {
	h.startShutdownMtx.Lock()
	defer h.startShutdownMtx.Unlock()

	h.logger.Debug("Starting Heartbeat routine")
	if h.cancel != nil {
		return alreadyRunningError
	}
	var ctx context.Context
	ctx, h.cancel = context.WithCancel(context.Background())
	go h.loop(ctx)
	return nil
}

func (h *Heartbeat) Shutdown() error {
	h.startShutdownMtx.Lock()
	defer h.startShutdownMtx.Unlock()

	h.logger.Debug("Stopping Heartbeat routine")
	if h.cancel == nil {
		return notRunningError
	}
	h.cancel()
	h.cancel = nil
	return nil
}

func (h *Heartbeat) loop(ctx context.Context) {
	tick := time.NewTicker(30 * time.Second)
	defer tick.Stop()

	// Start beat
	if err := h.generateHeartbeat(ctx); err != nil {
		h.logger.Error("Generating heartbeat failed", zap.Error(err))
	}

	for {
		select {
		case <-tick.C:
			if err := h.generateHeartbeat(ctx); err != nil {
				h.logger.Error("Generating heartbeat failed", zap.Error(err))
			}
		case <-ctx.Done():
			return
		}
	}

}

func (h *Heartbeat) generateHeartbeat(ctx context.Context) error {
	h.logger.Debug("Generating heartbeat")
	md := pmetric.NewMetrics()

	if err := h.addMetrics(ctx, md); err != nil {
		return err
	}

	return h.pushMetrics(ctx, md)
}