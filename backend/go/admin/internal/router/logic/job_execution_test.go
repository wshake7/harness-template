package logic

import (
	"admin/internal/mock"
	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	v1 "orm-crud/api/gen/go/pagination/v1"
	"orm-crud/gormc/mixin"
)

func mustMigrateJob(t *testing.T) *query.Query {
	t.Helper()
	return setupSQLiteDB(t,
		&models.JobExecution{},
		&models.JobSchedule{},
	)
}

func TestJobExecutionHandler_List_Empty(t *testing.T) {
	q := mustMigrateJob(t)
	query.SetDefault(q.JobExecution.UnderlyingDB())
	ctrl := gomock.NewController(t)
	h := NewJobExecutionHandler(q, mock.NewMockTemporalService(ctrl))

	result, err := h.List(newTestCtx(t), &v1.PagingRequest{})
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), result.Total)
}

func TestJobExecutionHandler_List_WithData(t *testing.T) {
	q := mustMigrateJob(t)
	query.SetDefault(q.JobExecution.UnderlyingDB())
	ctrl := gomock.NewController(t)
	h := NewJobExecutionHandler(q, mock.NewMockTemporalService(ctrl))

	q.JobExecution.Create(&models.JobExecution{
		AutoIncrementID: mixin.AutoIncrementID{ID: 1},
		JobCode:         "test-job",
		Status:          models.JobExecutionStatusRunning,
	})

	result, err := h.List(newTestCtx(t), &v1.PagingRequest{})
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), result.Total)
}

func TestJobExecutionHandler_Detail_NotFound(t *testing.T) {
	q := mustMigrateJob(t)
	query.SetDefault(q.JobExecution.UnderlyingDB())
	ctrl := gomock.NewController(t)
	h := NewJobExecutionHandler(q, mock.NewMockTemporalService(ctrl))

	_, err := h.Detail(newTestCtx(t), &ReqJobExecutionID{ID: 99})
	assert.Error(t, err)
}

func TestJobExecutionHandler_Cancel_NotFound(t *testing.T) {
	q := mustMigrateJob(t)
	query.SetDefault(q.JobExecution.UnderlyingDB())
	ctrl := gomock.NewController(t)
	h := NewJobExecutionHandler(q, mock.NewMockTemporalService(ctrl))

	err := h.Cancel(newTestCtx(t), &ReqJobExecutionID{ID: 99})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "执行记录不存在")
}

func TestJobExecutionHandler_Cancel_NotRunning(t *testing.T) {
	q := mustMigrateJob(t)
	query.SetDefault(q.JobExecution.UnderlyingDB())
	ctrl := gomock.NewController(t)
	h := NewJobExecutionHandler(q, mock.NewMockTemporalService(ctrl))

	q.JobExecution.Create(&models.JobExecution{
		AutoIncrementID: mixin.AutoIncrementID{ID: 1},
		JobCode:         "test-job",
		Status:          models.JobExecutionStatusSuccess,
	})

	err := h.Cancel(newTestCtx(t), &ReqJobExecutionID{ID: 1})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "仅运行中的执行记录可取消")
}

func TestJobExecutionHandler_Cancel_Success(t *testing.T) {
	q := mustMigrateJob(t)
	query.SetDefault(q.JobExecution.UnderlyingDB())
	ctrl := gomock.NewController(t)
	mockTemporal := mock.NewMockTemporalService(ctrl)
	h := NewJobExecutionHandler(q, mockTemporal)

	q.JobExecution.Create(&models.JobExecution{
		AutoIncrementID:     mixin.AutoIncrementID{ID: 1},
		JobCode:             "test-job",
		Status:              models.JobExecutionStatusRunning,
		TemporalWorkflowID:  "wf-1",
		TemporalRunID:       "run-1",
	})
	mockTemporal.EXPECT().CancelWorkflow(gomock.Any(), "wf-1", "run-1").Return(nil)

	err := h.Cancel(newTestCtx(t), &ReqJobExecutionID{ID: 1})
	assert.NoError(t, err)
}

func TestJobExecutionHandler_Retry_NotFound(t *testing.T) {
	q := mustMigrateJob(t)
	query.SetDefault(q.JobExecution.UnderlyingDB())
	ctrl := gomock.NewController(t)
	h := NewJobExecutionHandler(q, mock.NewMockTemporalService(ctrl))

	err := h.Retry(newTestCtx(t), &ReqJobExecutionID{ID: 99})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "执行记录不存在")
}

func TestJobExecutionHandler_Retry_NotRetryable(t *testing.T) {
	q := mustMigrateJob(t)
	query.SetDefault(q.JobExecution.UnderlyingDB())
	ctrl := gomock.NewController(t)
	h := NewJobExecutionHandler(q, mock.NewMockTemporalService(ctrl))

	q.JobExecution.Create(&models.JobExecution{
		AutoIncrementID: mixin.AutoIncrementID{ID: 1},
		JobCode:         "test-job",
		Status:          models.JobExecutionStatusRunning,
	})

	err := h.Retry(newTestCtx(t), &ReqJobExecutionID{ID: 1})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "仅失败、取消、超时的执行记录可重试")
}

func TestJobExecutionHandler_Retry_Success(t *testing.T) {
	q := mustMigrateJob(t)
	query.SetDefault(q.JobExecution.UnderlyingDB())
	ctrl := gomock.NewController(t)
	mockTemporal := mock.NewMockTemporalService(ctrl)
	h := NewJobExecutionHandler(q, mockTemporal)

	q.JobExecution.Create(&models.JobExecution{
		AutoIncrementID: mixin.AutoIncrementID{ID: 1},
		JobCode:         "test-job",
		Status:          models.JobExecutionStatusFailed,
		InputJSON:       []byte(`{"x":1}`),
		RetryCount:      1,
	})
	q.JobSchedule.Create(&models.JobSchedule{
		AutoIncrementID:           mixin.AutoIncrementID{ID: 1},
		JobCode:                   "test-job",
		JobName:                   "Test Job",
		WorkflowType:              "wf",
		TaskQueue:                 "admin",
		ScheduleType:              models.JobScheduleTypeCron,
		CronExpr:                  "* * * * *",
		Status:                    models.JobScheduleStatusEnabled,
		TemporalWorkflowIDPrefix:  "prefix",
	})
	mockTemporal.EXPECT().ExecuteWorkflow(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)

	err := h.Retry(newTestCtx(t), &ReqJobExecutionID{ID: 1})
	assert.NoError(t, err)
}

func TestJobExecutionHandler_Retry_InvalidJSON(t *testing.T) {
	q := mustMigrateJob(t)
	query.SetDefault(q.JobExecution.UnderlyingDB())
	ctrl := gomock.NewController(t)
	h := NewJobExecutionHandler(q, mock.NewMockTemporalService(ctrl))

	q.JobExecution.Create(&models.JobExecution{
		AutoIncrementID: mixin.AutoIncrementID{ID: 1},
		JobCode:         "test-job",
		Status:          models.JobExecutionStatusFailed,
		InputJSON:       []byte(`{bad`),
	})
	q.JobSchedule.Create(&models.JobSchedule{
		AutoIncrementID: mixin.AutoIncrementID{ID: 1},
		JobCode:         "test-job",
		JobName:         "Test Job",
		WorkflowType:    "wf",
		TaskQueue:       "admin",
		ScheduleType:    models.JobScheduleTypeCron,
		CronExpr:        "* * * * *",
		Status:          models.JobScheduleStatusEnabled,
	})

	err := h.Retry(newTestCtx(t), &ReqJobExecutionID{ID: 1})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "执行记录输入参数不是合法 JSON")
}

func TestIsRetryableJobExecution(t *testing.T) {
	assert.True(t, isRetryableJobExecution(models.JobExecutionStatusFailed))
	assert.True(t, isRetryableJobExecution(models.JobExecutionStatusCanceled))
	assert.True(t, isRetryableJobExecution(models.JobExecutionStatusTimeout))
	assert.False(t, isRetryableJobExecution(models.JobExecutionStatusRunning))
	assert.False(t, isRetryableJobExecution(models.JobExecutionStatusSuccess))
}
