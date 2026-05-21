package logic

import (
	"admin/internal/services/orm/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.temporal.io/api/serviceerror"
	"gorm.io/datatypes"
	v1 "orm-crud/api/gen/go/pagination/v1"
)

// --- normalize ---

func TestReqJobScheduleCreate_normalize(t *testing.T) {
	req := &ReqJobScheduleCreate{
		JobCode: " test ", JobName: " My Job ", WorkflowType: " wf ",
		TaskQueue: " q ", ScheduleType: " cron ", CronExpr: " * * * * * ",
		InputJSON: " {} ", Status: " enabled ",
		TemporalScheduleID: " id ", TemporalWorkflowIDPrefix: " pre ",
		Description: " desc ",
	}
	req.normalize()
	assert.Equal(t, "test", req.JobCode)
	assert.Equal(t, "My Job", req.JobName)
	assert.Equal(t, "wf", req.WorkflowType)
	assert.Equal(t, "q", req.TaskQueue)
	assert.Equal(t, "CRON", req.ScheduleType)
	assert.Equal(t, "* * * * *", req.CronExpr)
	assert.Equal(t, "{}", req.InputJSON)
	assert.Equal(t, "ENABLED", req.Status)
	assert.Equal(t, "id", req.TemporalScheduleID)
	assert.Equal(t, "pre", req.TemporalWorkflowIDPrefix)
	assert.Equal(t, "desc", req.Description)
}

func TestReqJobScheduleUpdate_normalize(t *testing.T) {
	name := " My Job "
	wf := " wf "
	queue := " q "
	sched := " cron "
	cron := " * * * * * "
	input := " {} "
	status := " enabled "
	prefix := " pre "
	desc := " desc "
	req := &ReqJobScheduleUpdate{
		JobName: &name, WorkflowType: &wf, TaskQueue: &queue,
		ScheduleType: &sched, CronExpr: &cron, InputJSON: &input,
		Status: &status, TemporalWorkflowIDPrefix: &prefix, Description: &desc,
	}
	req.normalize()
	assert.Equal(t, "My Job", *req.JobName)
	assert.Equal(t, "wf", *req.WorkflowType)
	assert.Equal(t, "q", *req.TaskQueue)
	assert.Equal(t, "CRON", *req.ScheduleType)
	assert.Equal(t, "* * * * *", *req.CronExpr)
	assert.Equal(t, "{}", *req.InputJSON)
	assert.Equal(t, "ENABLED", *req.Status)
	assert.Equal(t, "pre", *req.TemporalWorkflowIDPrefix)
	assert.Equal(t, "desc", *req.Description)
}

// --- validateJobSchedule ---

func TestValidateJobSchedule_Valid(t *testing.T) {
	err := validateJobSchedule(&models.JobSchedule{
		JobCode: "test", JobName: "Test", WorkflowType: "wf",
		TaskQueue: "q", ScheduleType: models.JobScheduleTypeCron,
		CronExpr: "* * * * *", Status: models.JobScheduleStatusEnabled,
	})
	assert.NoError(t, err)
}

func TestValidateJobSchedule_EmptyCode(t *testing.T) {
	err := validateJobSchedule(&models.JobSchedule{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "任务编码不能为空")
}

func TestValidateJobSchedule_EmptyName(t *testing.T) {
	err := validateJobSchedule(&models.JobSchedule{JobCode: "test"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "任务名称不能为空")
}

func TestValidateJobSchedule_EmptyWorkflowType(t *testing.T) {
	err := validateJobSchedule(&models.JobSchedule{JobCode: "t", JobName: "n"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Workflow 类型不能为空")
}

func TestValidateJobSchedule_InvalidScheduleType(t *testing.T) {
	err := validateJobSchedule(&models.JobSchedule{
		JobCode: "t", JobName: "n", WorkflowType: "w",
		TaskQueue: "q", ScheduleType: "INVALID", Status: models.JobScheduleStatusEnabled,
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "调度类型无效")
}

func TestValidateJobSchedule_InvalidStatus(t *testing.T) {
	err := validateJobSchedule(&models.JobSchedule{
		JobCode: "t", JobName: "n", WorkflowType: "w",
		TaskQueue: "q", ScheduleType: models.JobScheduleTypeCron,
		CronExpr: "* * * * *", Status: "INVALID",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "任务状态无效")
}

func TestValidateJobSchedule_CronWithoutExpr(t *testing.T) {
	err := validateJobSchedule(&models.JobSchedule{
		JobCode: "t", JobName: "n", WorkflowType: "w",
		TaskQueue: "q", ScheduleType: models.JobScheduleTypeCron,
		CronExpr: "", Status: models.JobScheduleStatusEnabled,
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "CRON 调度必须填写 cron 表达式")
}

func TestValidateJobSchedule_IntervalWithoutSeconds(t *testing.T) {
	err := validateJobSchedule(&models.JobSchedule{
		JobCode: "t", JobName: "n", WorkflowType: "w",
		TaskQueue: "q", ScheduleType: models.JobScheduleTypeInterval,
		Status: models.JobScheduleStatusEnabled,
	})
	assert.Error(t, err)
}

func TestValidateJobSchedule_OnceWithoutStartTime(t *testing.T) {
	err := validateJobSchedule(&models.JobSchedule{
		JobCode: "t", JobName: "n", WorkflowType: "w",
		TaskQueue: "q", ScheduleType: models.JobScheduleTypeOnce,
		Status: models.JobScheduleStatusEnabled,
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ONCE 调度必须填写开始时间")
}

func TestValidateJobSchedule_EndBeforeStart(t *testing.T) {
	now := time.Now()
	start := now
	end := now.Add(-1 * time.Hour)
	err := validateJobSchedule(&models.JobSchedule{
		JobCode: "t", JobName: "n", WorkflowType: "w",
		TaskQueue: "q", ScheduleType: models.JobScheduleTypeOnce,
		StartTime: &start, EndTime: &end,
		Status: models.JobScheduleStatusEnabled,
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "结束时间必须晚于开始时间")
}

// --- buildJobScheduleOptions ---

func TestBuildJobScheduleOptions(t *testing.T) {
	opts := buildJobScheduleOptions(map[string]string{
		"q1": "Queue 1",
		"q2": "Queue 2",
		"":   "should skip empty key",
	})
	assert.Len(t, opts, 2)
	assert.Equal(t, "Queue 1", opts[0].Label)
	assert.Equal(t, "q1", opts[0].Value)
}

func TestBuildJobScheduleOptions_EmptyLabel(t *testing.T) {
	opts := buildJobScheduleOptions(map[string]string{"test": ""})
	assert.Len(t, opts, 1)
	assert.Equal(t, "test", opts[0].Label)
}

// --- isTemporalNotFound ---

func TestIsTemporalNotFound(t *testing.T) {
	assert.True(t, isTemporalNotFound(&serviceerror.NotFound{}))
	assert.False(t, isTemporalNotFound(assert.AnError))
}

// --- normalizeJSONText / parseJSONValue / mustJSON ---

func TestNormalizeJSONText_Valid(t *testing.T) {
	j, v, err := normalizeJSONText(`{"key":"value"}`)
	assert.NoError(t, err)
	assert.Equal(t, datatypes.JSON(`{"key":"value"}`), j)
	assert.NotNil(t, v)
}

func TestNormalizeJSONText_Invalid(t *testing.T) {
	_, _, err := normalizeJSONText(`{bad json}`)
	assert.Error(t, err)
}

func TestNormalizeJSONText_Empty(t *testing.T) {
	j, v, err := normalizeJSONText("")
	assert.NoError(t, err)
	assert.Nil(t, j)
	assert.Nil(t, v)
}

func TestParseJSONValue(t *testing.T) {
	v, err := parseJSONValue(datatypes.JSON(`{"a":1}`))
	assert.NoError(t, err)
	assert.NotNil(t, v)
}

func TestParseJSONValue_Empty(t *testing.T) {
	v, err := parseJSONValue(nil)
	assert.NoError(t, err)
	assert.Nil(t, v)
}

func TestMustJSON(t *testing.T) {
	s := mustJSON(map[string]string{"k": "v"})
	assert.Contains(t, s, "k")
	assert.Contains(t, s, "v")
}

// --- buildScheduleSpec ---

func TestBuildScheduleSpec_Cron(t *testing.T) {
	spec, remaining, err := buildScheduleSpec(&models.JobSchedule{
		ScheduleType: models.JobScheduleTypeCron,
		CronExpr:     "0 * * * *",
	})
	assert.NoError(t, err)
	assert.Equal(t, 0, remaining)
	assert.Equal(t, []string{"0 * * * *"}, spec.CronExpressions)
}

func TestBuildScheduleSpec_Interval(t *testing.T) {
	ival := 300
	spec, remaining, err := buildScheduleSpec(&models.JobSchedule{
		ScheduleType:    models.JobScheduleTypeInterval,
		IntervalSeconds: &ival,
	})
	assert.NoError(t, err)
	assert.Equal(t, 0, remaining)
	assert.Len(t, spec.Intervals, 1)
}

func TestBuildScheduleSpec_Once(t *testing.T) {
	start := time.Date(2025, 1, 15, 10, 30, 45, 0, time.UTC)
	spec, remaining, err := buildScheduleSpec(&models.JobSchedule{
		ScheduleType: models.JobScheduleTypeOnce,
		StartTime:    &start,
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, remaining)
	assert.NotNil(t, spec.Calendars)
}

func TestBuildScheduleSpec_InvalidType(t *testing.T) {
	_, _, err := buildScheduleSpec(&models.JobSchedule{
		ScheduleType: "INVALID",
	})
	assert.Error(t, err)
}

// --- buildScheduleOptions ---

func TestBuildScheduleOptions(t *testing.T) {
	start := time.Now()
	opts, err := buildScheduleOptions(&models.JobSchedule{
		JobCode:         "test-code",
		JobName:         "Test",
		WorkflowType:    "wf",
		TaskQueue:       "q",
		ScheduleType:    models.JobScheduleTypeCron,
		CronExpr:        "* * * * *",
		TemporalScheduleID:       "test-sched",
		TemporalWorkflowIDPrefix: "test-prefix",
		Status:          models.JobScheduleStatusEnabled,
		StartTime:       &start,
	}, nil)
	assert.NoError(t, err)
	assert.Equal(t, "test-sched", opts.ID)
	assert.False(t, opts.Paused)
}

func TestBuildScheduleOptions_Paused(t *testing.T) {
	opts, err := buildScheduleOptions(&models.JobSchedule{
		JobCode:      "t",
		JobName:      "n",
		WorkflowType: "w",
		TaskQueue:    "q",
		ScheduleType: models.JobScheduleTypeCron,
		CronExpr:     "* * * * *",
		Status:       models.JobScheduleStatusDisabled,
	}, nil)
	assert.NoError(t, err)
	assert.True(t, opts.Paused)
}

// --- toModel ---

func TestReqJobScheduleCreate_ToModel(t *testing.T) {
	req := &ReqJobScheduleCreate{
		JobCode: "test", JobName: "Test", WorkflowType: "wf",
		TaskQueue: "q", ScheduleType: models.JobScheduleTypeCron,
		CronExpr: "* * * * *", InputJSON: `{"k":"v"}`,
		Status: models.JobScheduleStatusEnabled,
	}
	model, inputValue, err := req.toModel()
	assert.NoError(t, err)
	assert.Equal(t, "test", model.JobCode)
	assert.NotNil(t, inputValue)
}

func TestReqJobScheduleCreate_ToModel_Defaults(t *testing.T) {
	req := &ReqJobScheduleCreate{
		JobCode: "test", JobName: "Test", WorkflowType: "wf",
		TaskQueue: "q", ScheduleType: models.JobScheduleTypeCron,
		CronExpr: "* * * * *",
	}
	model, _, err := req.toModel()
	assert.NoError(t, err)
	assert.Equal(t, models.JobScheduleStatusEnabled, model.Status)
	assert.Equal(t, "test", model.TemporalScheduleID)
	assert.Equal(t, "test", model.TemporalWorkflowIDPrefix)
}

// --- appendPagingQuery ---

func TestAppendPagingQuery(t *testing.T) {
	req := &v1.PagingRequest{}
	appendPagingQuery(req, map[string]any{"status__not": "DELETED"})
	assert.NotNil(t, req.FilteringType)
	assert.Contains(t, req.GetQuery(), "DELETED")
}
