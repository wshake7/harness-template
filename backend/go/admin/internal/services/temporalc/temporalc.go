package temporalc

import (
	"admin/internal/config"
	"admin/internal/services/temporaljob"
	"fmt"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.uber.org/zap"
)

type Temporal struct {
	Client         client.Client
	ScheduleClient client.ScheduleClient
	Worker         worker.Worker
}

var Client *Temporal

func New(conf config.TemporalConfig, logger *zap.Logger) (*Temporal, error) {
	workflowClient, err := client.Dial(client.Options{
		HostPort:  conf.HostPort,
		Namespace: conf.Namespace,
		Identity:  conf.Identity,
		Logger:    zapTemporalLogger{logger: logger},
	})
	if err != nil {
		return nil, fmt.Errorf("dial temporal: %w", err)
	}

	Client = &Temporal{
		Client:         workflowClient,
		ScheduleClient: workflowClient.ScheduleClient(),
	}
	if conf.WorkerEnabled {
		Client.Worker = worker.New(workflowClient, conf.TaskQueue, worker.Options{})
		temporaljob.RegisterWorker(Client.Worker)
	}
	return Client, nil
}

func Close() {
	if Client == nil {
		return
	}
	if Client.Worker != nil {
		Client.Worker.Stop()
	}
	if Client.Client != nil {
		Client.Client.Close()
	}
	Client = nil
}

type zapTemporalLogger struct {
	logger *zap.Logger
}

func (l zapTemporalLogger) Debug(msg string, keyvals ...any) {
	l.zap().Debug(msg, temporalFields(keyvals)...)
}

func (l zapTemporalLogger) Info(msg string, keyvals ...any) {
	l.zap().Info(msg, temporalFields(keyvals)...)
}

func (l zapTemporalLogger) Warn(msg string, keyvals ...any) {
	l.zap().Warn(msg, temporalFields(keyvals)...)
}

func (l zapTemporalLogger) Error(msg string, keyvals ...any) {
	l.zap().Error(msg, temporalFields(keyvals)...)
}

func (l zapTemporalLogger) zap() *zap.Logger {
	if l.logger == nil {
		return zap.L().With(zap.String("module", "temporal"))
	}
	return l.logger.With(zap.String("module", "temporal"))
}

func temporalFields(keyvals []any) []zap.Field {
	fields := make([]zap.Field, 0, len(keyvals)/2+1)
	for i := 0; i < len(keyvals); i += 2 {
		key, ok := keyvals[i].(string)
		if !ok || key == "" {
			fields = append(fields, zap.Any("temporal_log_arg", keyvals[i]))
			continue
		}
		if i+1 >= len(keyvals) {
			fields = append(fields, zap.Any(key, nil))
			continue
		}
		fields = append(fields, zap.Any(key, keyvals[i+1]))
	}
	return fields
}
