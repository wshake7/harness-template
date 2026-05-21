package logic

import (
	"admin/internal/mock"
	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.temporal.io/api/serviceerror"
	v1 "orm-crud/api/gen/go/pagination/v1"
	"orm-crud/gormc/mixin"
)

func TestJobScheduleHandler_List_Empty(t *testing.T) {
	q := mustMigrateJob(t)
	query.SetDefault(q.JobSchedule.UnderlyingDB())
	ctrl := gomock.NewController(t)
	h := NewJobScheduleHandler(q, mock.NewMockTemporalService(ctrl))

	result, err := h.List(newTestCtx(t), &v1.PagingRequest{})
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), result.Total)
}

func TestJobScheduleHandler_List_WithData(t *testing.T) {
	q := mustMigrateJob(t)
	query.SetDefault(q.JobSchedule.UnderlyingDB())
	ctrl := gomock.NewController(t)
	h := NewJobScheduleHandler(q, mock.NewMockTemporalService(ctrl))

	q.JobSchedule.Create(&models.JobSchedule{
		AutoIncrementID: mixin.AutoIncrementID{ID: 1},
		JobCode:         "test-job",
		JobName:         "Test Job",
		WorkflowType:    "test",
		TaskQueue:       "admin",
		ScheduleType:    models.JobScheduleTypeCron,
		Status:          models.JobScheduleStatusEnabled,
	})

	result, err := h.List(newTestCtx(t), &v1.PagingRequest{})
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), result.Total)
}

func TestJobScheduleHandler_Detail_NotFound(t *testing.T) {
	q := mustMigrateJob(t)
	query.SetDefault(q.JobSchedule.UnderlyingDB())
	ctrl := gomock.NewController(t)
	h := NewJobScheduleHandler(q, mock.NewMockTemporalService(ctrl))

	_, err := h.Detail(newTestCtx(t), &ReqJobScheduleDetail{ID: 99})
	assert.Error(t, err)
}

func TestJobScheduleHandler_Detail_Success(t *testing.T) {
	q := mustMigrateJob(t)
	query.SetDefault(q.JobSchedule.UnderlyingDB())
	ctrl := gomock.NewController(t)
	h := NewJobScheduleHandler(q, mock.NewMockTemporalService(ctrl))

	q.JobSchedule.Create(&models.JobSchedule{
		AutoIncrementID: mixin.AutoIncrementID{ID: 1},
		JobCode:         "test-job",
		JobName:         "Test Job",
		WorkflowType:    "test",
		TaskQueue:       "admin",
		ScheduleType:    models.JobScheduleTypeCron,
		Status:          models.JobScheduleStatusEnabled,
	})

	result, err := h.Detail(newTestCtx(t), &ReqJobScheduleDetail{ID: 1})
	assert.NoError(t, err)
	assert.Equal(t, "test-job", result.JobCode)
}

func TestJobScheduleHandler_Options(t *testing.T) {
	q := mustMigrateJob(t)
	query.SetDefault(q.JobSchedule.UnderlyingDB())
	ctrl := gomock.NewController(t)
	h := NewJobScheduleHandler(q, mock.NewMockTemporalService(ctrl))

	result, err := h.Options(newTestCtx(t))
	assert.NoError(t, err)
	assert.NotEmpty(t, result.DefaultTaskQueue)
	assert.NotEmpty(t, result.WorkflowTypes)
	assert.NotEmpty(t, result.TaskQueues)
}

func TestJobScheduleHandler_Create_Invalid(t *testing.T) {
	q := mustMigrateJob(t)
	query.SetDefault(q.JobSchedule.UnderlyingDB())
	ctrl := gomock.NewController(t)
	h := NewJobScheduleHandler(q, mock.NewMockTemporalService(ctrl))

	err := h.Create(newTestCtx(t), &ReqJobScheduleCreate{})
	assert.Error(t, err)
}

func TestJobScheduleHandler_Create_Success(t *testing.T) {
	q := mustMigrateJob(t)
	query.SetDefault(q.JobSchedule.UnderlyingDB())
	ctrl := gomock.NewController(t)
	mockTemporal := mock.NewMockTemporalService(ctrl)
	h := NewJobScheduleHandler(q, mockTemporal)

	mockTemporal.EXPECT().UpdateSchedule(gomock.Any(), gomock.Any(), gomock.Any()).Return(assert.AnError)
	mockTemporal.EXPECT().CreateSchedule(gomock.Any(), gomock.Any()).Return(nil, nil)

	err := h.Create(newTestCtx(t), &ReqJobScheduleCreate{
		JobCode: "test-job", JobName: "Test Job", WorkflowType: "wf",
		TaskQueue: "admin", ScheduleType: models.JobScheduleTypeCron,
		CronExpr: "* * * * *",
	})
	assert.NoError(t, err)
}

func TestJobScheduleHandler_Del_NotFound(t *testing.T) {
	q := mustMigrateJob(t)
	query.SetDefault(q.JobSchedule.UnderlyingDB())
	ctrl := gomock.NewController(t)
	h := NewJobScheduleHandler(q, mock.NewMockTemporalService(ctrl))

	err := h.Del(newTestCtx(t), &ReqJobScheduleID{ID: 99})
	assert.Error(t, err)
}

func TestJobScheduleHandler_Switch_NotFound(t *testing.T) {
	q := mustMigrateJob(t)
	query.SetDefault(q.JobSchedule.UnderlyingDB())
	ctrl := gomock.NewController(t)
	h := NewJobScheduleHandler(q, mock.NewMockTemporalService(ctrl))

	err := h.Switch(newTestCtx(t), &ReqJobScheduleSwitch{ID: 99, Enabled: true})
	assert.Error(t, err)
}

func TestJobScheduleHandler_Sync_NotFound(t *testing.T) {
	q := mustMigrateJob(t)
	query.SetDefault(q.JobSchedule.UnderlyingDB())
	ctrl := gomock.NewController(t)
	h := NewJobScheduleHandler(q, mock.NewMockTemporalService(ctrl))

	err := h.Sync(newTestCtx(t), &ReqJobScheduleID{ID: 99})
	assert.Error(t, err)
}

func TestJobScheduleHandler_Trigger_NotFound(t *testing.T) {
	q := mustMigrateJob(t)
	query.SetDefault(q.JobSchedule.UnderlyingDB())
	ctrl := gomock.NewController(t)
	h := NewJobScheduleHandler(q, mock.NewMockTemporalService(ctrl))

	err := h.Trigger(newTestCtx(t), &ReqJobScheduleID{ID: 99})
	assert.Error(t, err)
}

func TestDeleteTemporalSchedule_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockTemporal := mock.NewMockTemporalService(ctrl)
	err := deleteTemporalSchedule(mockTemporal, nil, "")
	assert.NoError(t, err)
}

func TestSwitchTemporalSchedule(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockTemporal := mock.NewMockTemporalService(ctrl)
	mockTemporal.EXPECT().PauseSchedule(gomock.Any(), "sched-1", gomock.Any()).Return(nil)
	assert.NoError(t, switchTemporalSchedule(mockTemporal, nil, "sched-1", false))
}

func TestSwitchTemporalSchedule_Enable(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockTemporal := mock.NewMockTemporalService(ctrl)
	mockTemporal.EXPECT().UnpauseSchedule(gomock.Any(), "sched-1", gomock.Any()).Return(nil)
	assert.NoError(t, switchTemporalSchedule(mockTemporal, nil, "sched-1", true))
}

func TestTriggerTemporalSchedule(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockTemporal := mock.NewMockTemporalService(ctrl)
	mockTemporal.EXPECT().TriggerSchedule(gomock.Any(), "sched-1", gomock.Any()).Return(nil)
	assert.NoError(t, triggerTemporalSchedule(mockTemporal, nil, "sched-1"))
}

func TestFindActiveJobSchedule(t *testing.T) {
	q := mustMigrateJob(t)
	query.SetDefault(q.JobSchedule.UnderlyingDB())
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
	item, err := findActiveJobSchedule(q, 1)
	assert.NoError(t, err)
	assert.Equal(t, "test-job", item.JobCode)
}

func TestMergeCurrent(t *testing.T) {
	q := mustMigrateJob(t)
	query.SetDefault(q.JobSchedule.UnderlyingDB())
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
	name := "Updated Job"
	current, next, _, err := (&ReqJobScheduleUpdate{ID: 1, JobName: &name}).mergeCurrent(q)
	assert.NoError(t, err)
	assert.Equal(t, "Test Job", current.JobName)
	assert.Equal(t, "Updated Job", next.JobName)
}

func TestIsTemporalNotFound_WithTemporalError(t *testing.T) {
	assert.True(t, isTemporalNotFound(&serviceerror.NotFound{}))
}

func TestBuildScheduleSpec_WithWindow(t *testing.T) {
	start := time.Now()
	end := start.Add(time.Hour)
	spec, remaining, err := buildScheduleSpec(&models.JobSchedule{
		ScheduleType: models.JobScheduleTypeCron,
		CronExpr:     "* * * * *",
		StartTime:    &start,
		EndTime:      &end,
	})
	assert.NoError(t, err)
	assert.Equal(t, 0, remaining)
	assert.False(t, spec.StartAt.IsZero())
	assert.False(t, spec.EndAt.IsZero())
}

func TestJobScheduleHandler_Update_NotFound(t *testing.T) {
	q := mustMigrateJob(t)
	query.SetDefault(q.JobSchedule.UnderlyingDB())
	ctrl := gomock.NewController(t)
	h := NewJobScheduleHandler(q, mock.NewMockTemporalService(ctrl))
	name := "Updated"
	err := h.Update(newTestCtx(t), &ReqJobScheduleUpdate{ID: 99, JobName: &name})
	assert.Error(t, err)
}

func TestJobScheduleHandler_Del_Success(t *testing.T) {
	q := mustMigrateJob(t)
	query.SetDefault(q.JobSchedule.UnderlyingDB())
	ctrl := gomock.NewController(t)
	mockTemporal := mock.NewMockTemporalService(ctrl)
	h := NewJobScheduleHandler(q, mockTemporal)
	q.JobSchedule.Create(&models.JobSchedule{AutoIncrementID: mixin.AutoIncrementID{ID: 1}, JobCode: "test-job", JobName: "Test Job", WorkflowType: "wf", TaskQueue: "admin", ScheduleType: models.JobScheduleTypeCron, CronExpr: "* * * * *", Status: models.JobScheduleStatusEnabled, TemporalScheduleID: "sched-1"})
	mockTemporal.EXPECT().DeleteSchedule(gomock.Any(), "sched-1").Return(nil)
	err := h.Del(newTestCtx(t), &ReqJobScheduleID{ID: 1})
	assert.NoError(t, err)
}

func TestJobScheduleHandler_Sync_Success(t *testing.T) {
	q := mustMigrateJob(t)
	query.SetDefault(q.JobSchedule.UnderlyingDB())
	ctrl := gomock.NewController(t)
	mockTemporal := mock.NewMockTemporalService(ctrl)
	h := NewJobScheduleHandler(q, mockTemporal)
	q.JobSchedule.Create(&models.JobSchedule{AutoIncrementID: mixin.AutoIncrementID{ID: 1}, JobCode: "test-job", JobName: "Test Job", WorkflowType: "wf", TaskQueue: "admin", ScheduleType: models.JobScheduleTypeCron, CronExpr: "* * * * *", Status: models.JobScheduleStatusEnabled, TemporalScheduleID: "sched-1", TemporalWorkflowIDPrefix: "prefix", InputJSON: []byte(`{"a":1}`)})
	mockTemporal.EXPECT().UpdateSchedule(gomock.Any(), "sched-1", gomock.Any()).Return(assert.AnError)
	mockTemporal.EXPECT().CreateSchedule(gomock.Any(), gomock.Any()).Return(nil, nil)
	err := h.Sync(newTestCtx(t), &ReqJobScheduleID{ID: 1})
	assert.NoError(t, err)
}

func TestJobScheduleHandler_Trigger_Success(t *testing.T) {
	q := mustMigrateJob(t)
	query.SetDefault(q.JobSchedule.UnderlyingDB())
	ctrl := gomock.NewController(t)
	mockTemporal := mock.NewMockTemporalService(ctrl)
	h := NewJobScheduleHandler(q, mockTemporal)
	q.JobSchedule.Create(&models.JobSchedule{AutoIncrementID: mixin.AutoIncrementID{ID: 1}, JobCode: "test-job", JobName: "Test Job", WorkflowType: "wf", TaskQueue: "admin", ScheduleType: models.JobScheduleTypeCron, CronExpr: "* * * * *", Status: models.JobScheduleStatusEnabled, TemporalScheduleID: "sched-1"})
	mockTemporal.EXPECT().TriggerSchedule(gomock.Any(), "sched-1", gomock.Any()).Return(nil)
	err := h.Trigger(newTestCtx(t), &ReqJobScheduleID{ID: 1})
	assert.NoError(t, err)
}

func TestJobScheduleHandler_Update_InvalidMergedModel(t *testing.T) {
	q := mustMigrateJob(t)
	query.SetDefault(q.JobSchedule.UnderlyingDB())
	ctrl := gomock.NewController(t)
	h := NewJobScheduleHandler(q, mock.NewMockTemporalService(ctrl))

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

	scheduleType := models.JobScheduleTypeInterval
	interval := 0
	err := h.Update(newTestCtx(t), &ReqJobScheduleUpdate{ID: 1, ScheduleType: &scheduleType, IntervalSeconds: &interval})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "INTERVAL 调度必须填写大于 0 的间隔秒数")
}

func TestJobScheduleHandler_Update_Success(t *testing.T) {
	q := mustMigrateJob(t)
	query.SetDefault(q.JobSchedule.UnderlyingDB())
	ctrl := gomock.NewController(t)
	mockTemporal := mock.NewMockTemporalService(ctrl)
	h := NewJobScheduleHandler(q, mockTemporal)

	q.JobSchedule.Create(&models.JobSchedule{
		AutoIncrementID:          mixin.AutoIncrementID{ID: 1},
		JobCode:                  "test-job",
		JobName:                  "Test Job",
		WorkflowType:             "wf",
		TaskQueue:                "admin",
		ScheduleType:             models.JobScheduleTypeCron,
		CronExpr:                 "* * * * *",
		Status:                   models.JobScheduleStatusEnabled,
		TemporalScheduleID:       "sched-1",
		TemporalWorkflowIDPrefix: "prefix",
		InputJSON:                []byte(`{"a":1}`),
	})
	mockTemporal.EXPECT().UpdateSchedule(gomock.Any(), "sched-1", gomock.Any()).Return(nil)

	name := "Updated Job"
	err := h.Update(newTestCtx(t), &ReqJobScheduleUpdate{ID: 1, JobName: &name})
	assert.NoError(t, err)
}

func TestJobScheduleHandler_Switch_EnableSuccess(t *testing.T) {
	q := mustMigrateJob(t)
	query.SetDefault(q.JobSchedule.UnderlyingDB())
	ctrl := gomock.NewController(t)
	mockTemporal := mock.NewMockTemporalService(ctrl)
	h := NewJobScheduleHandler(q, mockTemporal)

	q.JobSchedule.Create(&models.JobSchedule{
		AutoIncrementID:    mixin.AutoIncrementID{ID: 1},
		JobCode:            "test-job",
		JobName:            "Test Job",
		WorkflowType:       "wf",
		TaskQueue:          "admin",
		ScheduleType:       models.JobScheduleTypeCron,
		CronExpr:           "* * * * *",
		Status:             models.JobScheduleStatusDisabled,
		TemporalScheduleID: "sched-1",
	})
	mockTemporal.EXPECT().UnpauseSchedule(gomock.Any(), "sched-1", gomock.Any()).Return(nil)

	err := h.Switch(newTestCtx(t), &ReqJobScheduleSwitch{ID: 1, Enabled: true})
	assert.NoError(t, err)
}

func TestJobScheduleHandler_Switch_NotFoundRecoverySuccess(t *testing.T) {
	q := mustMigrateJob(t)
	query.SetDefault(q.JobSchedule.UnderlyingDB())
	ctrl := gomock.NewController(t)
	mockTemporal := mock.NewMockTemporalService(ctrl)
	h := NewJobScheduleHandler(q, mockTemporal)

	q.JobSchedule.Create(&models.JobSchedule{
		AutoIncrementID:          mixin.AutoIncrementID{ID: 1},
		JobCode:                  "test-job",
		JobName:                  "Test Job",
		WorkflowType:             "wf",
		TaskQueue:                "admin",
		ScheduleType:             models.JobScheduleTypeCron,
		CronExpr:                 "* * * * *",
		Status:                   models.JobScheduleStatusEnabled,
		TemporalScheduleID:       "sched-1",
		TemporalWorkflowIDPrefix: "prefix",
		InputJSON:                []byte(`{"a":1}`),
	})
	mockTemporal.EXPECT().UnpauseSchedule(gomock.Any(), "sched-1", gomock.Any()).Return(&serviceerror.NotFound{})
	mockTemporal.EXPECT().UpdateSchedule(gomock.Any(), "sched-1", gomock.Any()).Return(assert.AnError)
	mockTemporal.EXPECT().CreateSchedule(gomock.Any(), gomock.Any()).Return(nil, nil)
	mockTemporal.EXPECT().UnpauseSchedule(gomock.Any(), "sched-1", gomock.Any()).Return(nil)

	err := h.Switch(newTestCtx(t), &ReqJobScheduleSwitch{ID: 1, Enabled: true})
	assert.NoError(t, err)
}

func TestJobScheduleHandler_Switch_NotFoundRecoveryBadJSON(t *testing.T) {
	q := mustMigrateJob(t)
	query.SetDefault(q.JobSchedule.UnderlyingDB())
	ctrl := gomock.NewController(t)
	mockTemporal := mock.NewMockTemporalService(ctrl)
	h := NewJobScheduleHandler(q, mockTemporal)

	q.JobSchedule.Create(&models.JobSchedule{
		AutoIncrementID:    mixin.AutoIncrementID{ID: 1},
		JobCode:            "test-job",
		JobName:            "Test Job",
		WorkflowType:       "wf",
		TaskQueue:          "admin",
		ScheduleType:       models.JobScheduleTypeCron,
		CronExpr:           "* * * * *",
		Status:             models.JobScheduleStatusEnabled,
		TemporalScheduleID: "sched-1",
		InputJSON:          []byte(`{bad`),
	})
	mockTemporal.EXPECT().UnpauseSchedule(gomock.Any(), "sched-1", gomock.Any()).Return(&serviceerror.NotFound{})

	err := h.Switch(newTestCtx(t), &ReqJobScheduleSwitch{ID: 1, Enabled: true})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "输入参数不是合法 JSON")
}
