package services

import (
	"admin/internal/config"
	"admin/internal/services/temporalc"
	"context"
	"fmt"

	"go.temporal.io/sdk/client"
	"go.uber.org/zap"
)

type Temporal struct {
	conf   config.TemporalConfig
	client *temporalc.Temporal
}

func NewTemporal(conf config.TemporalConfig) *Temporal {
	return &Temporal{conf: conf}
}

func (t *Temporal) Start(ctx context.Context) error {
	if !t.conf.Enabled {
		return nil
	}
	temporalClient, err := temporalc.New(t.conf, zap.L())
	if err != nil {
		return err
	}
	if temporalClient.Worker != nil {
		if err = temporalClient.Worker.Start(); err != nil {
			temporalc.Close()
			return fmt.Errorf("start temporal worker: %w", err)
		}
	}
	t.client = temporalClient
	return nil
}

func (t *Temporal) String() string {
	return "temporal"
}

func (t *Temporal) State(ctx context.Context) (string, error) {
	if !t.conf.Enabled {
		return "DISABLED", nil
	}
	if t.client == nil || t.client.Client == nil {
		return "UNHEALTHY", fmt.Errorf("temporal client not initialized")
	}
	if _, err := t.client.Client.CheckHealth(ctx, &client.CheckHealthRequest{}); err != nil {
		return "UNHEALTHY", fmt.Errorf("temporal health check failed: %w", err)
	}
	return "HEALTHY", nil
}

func (t *Temporal) Terminate(ctx context.Context) error {
	if !t.conf.Enabled {
		return nil
	}
	temporalc.Close()
	return nil
}
