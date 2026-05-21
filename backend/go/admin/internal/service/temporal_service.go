package service

import (
	"context"

	"admin/internal/fiberc/res"
	"admin/internal/services/temporalc"

	"go.temporal.io/sdk/client"
)

//go:generate mockgen -source=temporal_service.go -destination=../mock/mock_temporal_service.go -package=mock -typed

type TemporalService interface {
	CancelWorkflow(ctx context.Context, workflowID, runID string) error
	ExecuteWorkflow(ctx context.Context, options client.StartWorkflowOptions, workflow any, args ...any) (client.WorkflowRun, error)
	CreateSchedule(ctx context.Context, options client.ScheduleOptions) (client.ScheduleHandle, error)
	UpdateSchedule(ctx context.Context, scheduleID string, options client.ScheduleUpdateOptions) error
	DeleteSchedule(ctx context.Context, scheduleID string) error
	PauseSchedule(ctx context.Context, scheduleID string, options client.SchedulePauseOptions) error
	UnpauseSchedule(ctx context.Context, scheduleID string, options client.ScheduleUnpauseOptions) error
	TriggerSchedule(ctx context.Context, scheduleID string, options client.ScheduleTriggerOptions) error
}

type temporalServiceImpl struct{}

func NewTemporalService() TemporalService {
	return &temporalServiceImpl{}
}

func requireTemporalService() (*temporalc.Temporal, error) {
	if temporalc.Client == nil || temporalc.Client.Client == nil || temporalc.Client.ScheduleClient == nil {
		return nil, res.FailMsg("Temporal 服务未启用")
	}
	return temporalc.Client, nil
}

func (s *temporalServiceImpl) CancelWorkflow(ctx context.Context, workflowID, runID string) error {
	service, err := requireTemporalService()
	if err != nil {
		return err
	}
	return service.Client.CancelWorkflow(ctx, workflowID, runID)
}

func (s *temporalServiceImpl) ExecuteWorkflow(ctx context.Context, options client.StartWorkflowOptions, workflow any, args ...any) (client.WorkflowRun, error) {
	service, err := requireTemporalService()
	if err != nil {
		return nil, err
	}
	return service.Client.ExecuteWorkflow(ctx, options, workflow, args...)
}

func (s *temporalServiceImpl) CreateSchedule(ctx context.Context, options client.ScheduleOptions) (client.ScheduleHandle, error) {
	service, err := requireTemporalService()
	if err != nil {
		return nil, err
	}
	return service.ScheduleClient.Create(ctx, options)
}

func (s *temporalServiceImpl) UpdateSchedule(ctx context.Context, scheduleID string, options client.ScheduleUpdateOptions) error {
	service, err := requireTemporalService()
	if err != nil {
		return err
	}
	return service.ScheduleClient.GetHandle(ctx, scheduleID).Update(ctx, options)
}

func (s *temporalServiceImpl) DeleteSchedule(ctx context.Context, scheduleID string) error {
	if scheduleID == "" {
		return nil
	}
	service, err := requireTemporalService()
	if err != nil {
		return err
	}
	return service.ScheduleClient.GetHandle(ctx, scheduleID).Delete(ctx)
}

func (s *temporalServiceImpl) PauseSchedule(ctx context.Context, scheduleID string, options client.SchedulePauseOptions) error {
	service, err := requireTemporalService()
	if err != nil {
		return err
	}
	return service.ScheduleClient.GetHandle(ctx, scheduleID).Pause(ctx, options)
}

func (s *temporalServiceImpl) UnpauseSchedule(ctx context.Context, scheduleID string, options client.ScheduleUnpauseOptions) error {
	service, err := requireTemporalService()
	if err != nil {
		return err
	}
	return service.ScheduleClient.GetHandle(ctx, scheduleID).Unpause(ctx, options)
}

func (s *temporalServiceImpl) TriggerSchedule(ctx context.Context, scheduleID string, options client.ScheduleTriggerOptions) error {
	service, err := requireTemporalService()
	if err != nil {
		return err
	}
	return service.ScheduleClient.GetHandle(ctx, scheduleID).Trigger(ctx, options)
}
