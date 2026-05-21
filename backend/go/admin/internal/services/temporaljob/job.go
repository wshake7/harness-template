package temporaljob

import (
	"context"
	"fmt"
	"time"

	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"

	"github.com/bytedance/sonic"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"gorm.io/datatypes"
)

const DispatchWorkflowName = "JobDispatchWorkflow"

func WorkflowTypeOptions() map[string]string {
	return map[string]string{
		PrintCountWorkflowName: PrintCountWorkflowName,
	}
}

type DispatchInput struct {
	JobCode          string `json:"jobCode"`
	WorkflowType     string `json:"workflowType"`
	TaskQueue        string `json:"taskQueue"`
	WorkflowIDPrefix string `json:"workflowIDPrefix"`
	Input            any    `json:"input"`
	RetryCount       int    `json:"retryCount"`
}

type CreateExecutionInput struct {
	JobCode            string    `json:"jobCode"`
	TemporalWorkflowID string    `json:"temporalWorkflowID"`
	TemporalRunID      string    `json:"temporalRunID"`
	TriggerTime        time.Time `json:"triggerTime"`
	Input              any       `json:"input"`
	RetryCount         int       `json:"retryCount"`
}

type CompleteExecutionInput struct {
	ID           uint64 `json:"id"`
	Status       string `json:"status"`
	Result       any    `json:"result"`
	ErrorMessage string `json:"errorMessage"`
}

type CreateExecutionResult struct {
	ID uint64 `json:"id"`
}

func RegisterWorker(w workerRegistry) {
	w.RegisterWorkflowWithOptions(DispatchWorkflow, workflow.RegisterOptions{Name: DispatchWorkflowName})
	w.RegisterWorkflowWithOptions(PrintCountWorkflow, workflow.RegisterOptions{Name: PrintCountWorkflowName})
	w.RegisterActivity(CreateExecution)
	w.RegisterActivity(CompleteExecution)
}

type workerRegistry interface {
	RegisterWorkflowWithOptions(workflowFunc any, options workflow.RegisterOptions)
	RegisterActivity(activityFunc any)
}

func DispatchWorkflow(ctx workflow.Context, input DispatchInput) error {
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2,
			MaximumAttempts:    3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	workflowID := buildWorkflowID(input, workflow.Now(ctx))
	childOptions := workflow.ChildWorkflowOptions{
		WorkflowID: workflowID,
		TaskQueue:  input.TaskQueue,
	}
	childCtx := workflow.WithChildOptions(ctx, childOptions)

	args := []any{}
	if input.Input != nil {
		args = append(args, input.Input)
	}
	child := workflow.ExecuteChildWorkflow(childCtx, input.WorkflowType, args...)

	var execution workflow.Execution
	if err := child.GetChildWorkflowExecution().Get(ctx, &execution); err != nil {
		return err
	}

	createInput := CreateExecutionInput{
		JobCode:            input.JobCode,
		TemporalWorkflowID: execution.ID,
		TemporalRunID:      execution.RunID,
		TriggerTime:        workflow.Now(ctx),
		Input:              input.Input,
		RetryCount:         input.RetryCount,
	}
	var createResult CreateExecutionResult
	if err := workflow.ExecuteActivity(ctx, CreateExecution, createInput).Get(ctx, &createResult); err != nil {
		return err
	}

	var childResult any
	err := child.Get(ctx, &childResult)
	completeInput := CompleteExecutionInput{
		ID:     createResult.ID,
		Result: childResult,
	}
	if err != nil {
		completeInput.Status = executionErrorStatus(err)
		completeInput.ErrorMessage = err.Error()
		_ = workflow.ExecuteActivity(ctx, CompleteExecution, completeInput).Get(ctx, nil)
		return err
	}

	completeInput.Status = models.JobExecutionStatusSuccess
	return workflow.ExecuteActivity(ctx, CompleteExecution, completeInput).Get(ctx, nil)
}

func CreateExecution(ctx context.Context, input CreateExecutionInput) (*CreateExecutionResult, error) {
	inputJSON, err := marshalNullableJSON(input.Input)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	execution := &models.JobExecution{
		JobCode:            input.JobCode,
		TemporalWorkflowID: input.TemporalWorkflowID,
		TemporalRunID:      input.TemporalRunID,
		TriggerTime:        input.TriggerTime,
		StartTime:          &now,
		Status:             models.JobExecutionStatusRunning,
		InputJSON:          inputJSON,
		RetryCount:         input.RetryCount,
	}
	if err = query.JobExecution.WithContext(ctx).Create(execution); err != nil {
		return nil, err
	}
	return &CreateExecutionResult{ID: execution.ID}, nil
}

func CompleteExecution(ctx context.Context, input CompleteExecutionInput) error {
	resultJSON, err := marshalNullableJSON(input.Result)
	if err != nil {
		return err
	}
	jobExecution := query.JobExecution
	_, err = jobExecution.WithContext(ctx).
		Where(jobExecution.ID.Eq(input.ID)).
		UpdateSimple(
			jobExecution.Status.Value(input.Status),
			jobExecution.ResultJSON.Value(resultJSON),
			jobExecution.ErrorMessage.Value(input.ErrorMessage),
			jobExecution.EndTime.Value(time.Now()),
		)
	return err
}

func buildWorkflowID(input DispatchInput, now time.Time) string {
	prefix := input.WorkflowIDPrefix
	if prefix == "" {
		prefix = input.JobCode
	}
	return fmt.Sprintf("%s-%d", prefix, now.UnixNano())
}

func executionErrorStatus(err error) string {
	switch {
	case temporal.IsCanceledError(err):
		return models.JobExecutionStatusCanceled
	case temporal.IsTimeoutError(err):
		return models.JobExecutionStatusTimeout
	default:
		return models.JobExecutionStatusFailed
	}
}

func marshalNullableJSON(value any) (datatypes.JSON, error) {
	if value == nil {
		return nil, nil
	}
	data, err := sonic.Marshal(value)
	if err != nil {
		return nil, err
	}
	return data, nil
}
