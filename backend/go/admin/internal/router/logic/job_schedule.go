package logic

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"admin/internal/config"
	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/res"
	"admin/internal/service"
	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"
	"admin/internal/services/temporaljob"
	"github.com/bytedance/sonic"
	"go-common/utils/str"
	"go.temporal.io/api/serviceerror"
	"go.temporal.io/sdk/client"
	"go.uber.org/zap"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	v1 "orm-crud/api/gen/go/pagination/v1"
	"orm-crud/gormc"
)

type JobScheduleHandler struct {
	Q        *query.Query
	Temporal service.TemporalService
}

func NewJobScheduleHandler(q *query.Query, temporal service.TemporalService) *JobScheduleHandler {
	return &JobScheduleHandler{Q: q, Temporal: temporal}
}

type RespJobSchedule struct {
	models.JobSchedule
	CanWrite  bool `json:"canWrite"`
	CanDelete bool `json:"canDelete"`
}

type RespJobScheduleOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type RespJobScheduleOptions struct {
	WorkflowTypes    []RespJobScheduleOption `json:"workflowTypes"`
	TaskQueues       []RespJobScheduleOption `json:"taskQueues"`
	DefaultTaskQueue string                  `json:"defaultTaskQueue"`
}

type ReqJobScheduleDetail struct {
	ID uint64 `json:"id" binding:"required" binding_msg:"required=请求错误"`
}

type ReqJobScheduleCreate struct {
	JobCode                  string     `json:"jobCode" change:"任务编码" binding:"required,max=128" binding_msg:"required=任务编码不能为空,max=任务编码最多128位"`
	JobName                  string     `json:"jobName" change:"任务名称" binding:"required,max=255" binding_msg:"required=任务名称不能为空,max=任务名称最多255位"`
	WorkflowType             string     `json:"workflowType" change:"Workflow类型" binding:"required,max=255" binding_msg:"required=Workflow 类型不能为空,max=Workflow 类型最多255位"`
	TaskQueue                string     `json:"taskQueue" change:"Task Queue" binding:"required,max=255" binding_msg:"required=Task Queue 不能为空,max=Task Queue 最多255位"`
	ScheduleType             string     `json:"scheduleType" change:"调度类型" binding:"required,max=32" binding_msg:"required=调度类型不能为空,max=调度类型最多32位"`
	CronExpr                 string     `json:"cronExpr" change:"Cron表达式" binding:"max=128" binding_msg:"max=cron 表达式最多128位"`
	IntervalSeconds          *int       `json:"intervalSeconds" change:"间隔秒数"`
	StartTime                *time.Time `json:"startTime" change:"开始时间"`
	EndTime                  *time.Time `json:"endTime" change:"结束时间"`
	InputJSON                string     `json:"inputJSON" change:"输入参数"`
	Status                   string     `json:"status" change:"状态" binding:"max=32" binding_msg:"max=状态最多32位"`
	TemporalScheduleID       string     `json:"temporalScheduleID" change:"Temporal Schedule ID" binding:"max=255" binding_msg:"max=Temporal Schedule ID 最多255位"`
	TemporalWorkflowIDPrefix string     `json:"temporalWorkflowIDPrefix" change:"Workflow ID前缀" binding:"max=255" binding_msg:"max=Workflow ID 前缀最多255位"`
	Description              string     `json:"description" change:"描述" binding:"max=512" binding_msg:"max=描述最多512位"`
}

type ReqJobScheduleUpdate struct {
	ID                       uint64     `json:"id" binding:"required" binding_msg:"required=请求错误"`
	JobName                  *string    `json:"jobName" change:"任务名称" binding:"omitempty,max=255" binding_msg:"max=任务名称最多255位"`
	WorkflowType             *string    `json:"workflowType" change:"Workflow类型" binding:"omitempty,max=255" binding_msg:"max=Workflow 类型最多255位"`
	TaskQueue                *string    `json:"taskQueue" change:"Task Queue" binding:"omitempty,max=255" binding_msg:"max=Task Queue 最多255位"`
	ScheduleType             *string    `json:"scheduleType" change:"调度类型" binding:"omitempty,max=32" binding_msg:"max=调度类型最多32位"`
	CronExpr                 *string    `json:"cronExpr" change:"Cron表达式" binding:"omitempty,max=128" binding_msg:"max=cron 表达式最多128位"`
	IntervalSeconds          *int       `json:"intervalSeconds" change:"间隔秒数"`
	StartTime                *time.Time `json:"startTime" change:"开始时间"`
	EndTime                  *time.Time `json:"endTime" change:"结束时间"`
	InputJSON                *string    `json:"inputJSON" change:"输入参数"`
	Status                   *string    `json:"status" change:"状态" binding:"omitempty,max=32" binding_msg:"max=状态最多32位"`
	TemporalWorkflowIDPrefix *string    `json:"temporalWorkflowIDPrefix" change:"Workflow ID前缀" binding:"omitempty,max=255" binding_msg:"max=Workflow ID 前缀最多255位"`
	Description              *string    `json:"description" change:"描述" binding:"omitempty,max=512" binding_msg:"max=描述最多512位"`
}

type ReqJobScheduleID struct {
	ID uint64 `json:"id" binding:"required" binding_msg:"required=请求错误"`
}

type ReqJobScheduleSwitch struct {
	ID      uint64 `json:"id" binding:"required" binding_msg:"required=请求错误"`
	Enabled bool   `json:"enabled" change:"启用状态"`
}

// @Summary 获取任务调度表单选项
// @Tags JobSchedule
// @Produce json
// @Success 200 {object} res.Response{data=RespJobScheduleOptions} "成功"
// @Router /api/sys/job/schedule/options [post]
func (h *JobScheduleHandler) Options(ctx *handler.Ctx) (*RespJobScheduleOptions, error) {
	defaultTaskQueue := strings.TrimSpace(config.Conf.Temporal.TaskQueue)
	if defaultTaskQueue == "" {
		defaultTaskQueue = "admin"
	}
	taskQueues := map[string]string{
		"admin":          "admin",
		defaultTaskQueue: defaultTaskQueue,
	}
	return &RespJobScheduleOptions{
		WorkflowTypes:    buildJobScheduleOptions(temporaljob.WorkflowTypeOptions()),
		TaskQueues:       buildJobScheduleOptions(taskQueues),
		DefaultTaskQueue: defaultTaskQueue,
	}, nil
}

// @Summary 获取任务调度配置分页列表
// @Tags JobSchedule
// @Accept json
// @Produce json
// @Param req body v1.PagingRequest true "分页参数"
// @Success 200 {object} res.Response{data=gormc.PagingResult[RespJobSchedule]} "成功"
// @Router /api/sys/job/schedule/list [post]
func (h *JobScheduleHandler) List(ctx *handler.Ctx, req *v1.PagingRequest) (*gormc.PagingResult[RespJobSchedule], error) {
	if req.GetOrderBy() == "" {
		orderBy := "id desc"
		req.OrderBy = &orderBy
	}
	appendPagingQuery(req, map[string]any{"status__not": models.JobScheduleStatusDeleted})

	pagination, err := h.Q.JobSchedule.PageWithPaging(req)
	if err != nil {
		ctx.L().Error("query job schedule list fail", zap.Error(err))
		return nil, res.FailDefault
	}

	items := make([]*RespJobSchedule, 0, len(pagination.Items))
	for _, item := range pagination.Items {
		if item == nil {
			continue
		}
		items = append(items, &RespJobSchedule{
			JobSchedule: *item,
			CanWrite:    true,
			CanDelete:   true,
		})
	}
	return &gormc.PagingResult[RespJobSchedule]{Items: items, Total: pagination.Total}, nil
}

// @Summary 获取任务调度配置详情
// @Tags JobSchedule
// @Accept json
// @Produce json
// @Param req body ReqJobScheduleDetail true "任务ID"
// @Success 200 {object} res.Response{data=models.JobSchedule} "成功"
// @Router /api/sys/job/schedule/detail [post]
func (h *JobScheduleHandler) Detail(ctx *handler.Ctx, req *ReqJobScheduleDetail) (*models.JobSchedule, error) {
	jobSchedule := h.Q.JobSchedule
	item, err := jobSchedule.Where(jobSchedule.ID.Eq(req.ID), jobSchedule.Status.Neq(models.JobScheduleStatusDeleted)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, res.FailMsg("任务不存在")
		}
		ctx.L().Error("query job schedule detail fail", zap.Error(err), zap.Uint64("id", req.ID))
		return nil, res.FailDefault
	}
	return item, nil
}

// @Summary 创建任务调度配置
// @Tags JobSchedule
// @Accept json
// @Produce json
// @Param req body ReqJobScheduleCreate true "创建参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/sys/job/schedule/create [post]
func (h *JobScheduleHandler) Create(ctx *handler.Ctx, req *ReqJobScheduleCreate) error {
	req.normalize()
	model, inputValue, err := req.toModel()
	if err != nil {
		return err
	}
	if err = validateJobSchedule(model); err != nil {
		return err
	}
	if err = syncTemporalSchedule(h.Temporal, ctx.Context(), model, inputValue, false); err != nil {
		ctx.L().Error("create temporal schedule fail", zap.Error(err), zap.String("jobCode", model.JobCode))
		return res.FailMsg("同步 Temporal Schedule 失败")
	}

	err = h.Q.JobSchedule.Create(model)
	if err != nil {
		_ = deleteTemporalSchedule(h.Temporal, ctx.Context(), model.TemporalScheduleID)
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return res.FailMsg("任务编码已存在")
		}
		return res.FailDefault
	}
	return nil
}

// @Summary 更新任务调度配置
// @Tags JobSchedule
// @Accept json
// @Produce json
// @Param req body ReqJobScheduleUpdate true "更新参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/sys/job/schedule/update [post]
func (h *JobScheduleHandler) Update(ctx *handler.Ctx, req *ReqJobScheduleUpdate) error {
	req.normalize()
	current, next, inputValue, err := req.mergeCurrent(h.Q)
	if err != nil {
		return err
	}
	if err = validateJobSchedule(next); err != nil {
		return err
	}
	if err = syncTemporalSchedule(h.Temporal, ctx.Context(), next, inputValue, current.TemporalScheduleID == ""); err != nil {
		ctx.L().Error("update temporal schedule fail", zap.Error(err), zap.Uint64("id", req.ID))
		return res.FailMsg("同步 Temporal Schedule 失败")
	}

	jobSchedule := h.Q.JobSchedule
	info, err := jobSchedule.Where(jobSchedule.ID.Eq(req.ID), jobSchedule.Status.Neq(models.JobScheduleStatusDeleted)).Updates(map[string]any{
		"job_name":                    next.JobName,
		"workflow_type":               next.WorkflowType,
		"task_queue":                  next.TaskQueue,
		"schedule_type":               next.ScheduleType,
		"cron_expr":                   next.CronExpr,
		"interval_seconds":            next.IntervalSeconds,
		"start_time":                  next.StartTime,
		"end_time":                    next.EndTime,
		"input_json":                  next.InputJSON,
		"status":                      next.Status,
		"temporal_schedule_id":        next.TemporalScheduleID,
		"temporal_workflow_id_prefix": next.TemporalWorkflowIDPrefix,
		"description":                 next.Description,
	})
	if err != nil {
		return res.FailDefault
	}
	if info.RowsAffected == 0 {
		return res.FailMsg("任务不存在")
	}
	return nil
}

// @Summary 删除任务调度配置
// @Tags JobSchedule
// @Accept json
// @Produce json
// @Param req body ReqJobScheduleID true "任务ID"
// @Success 200 {object} res.Response "成功"
// @Router /api/sys/job/schedule/del [post]
func (h *JobScheduleHandler) Del(ctx *handler.Ctx, req *ReqJobScheduleID) error {
	current, err := findActiveJobSchedule(h.Q, req.ID)
	if err != nil {
		return err
	}
	if err = deleteTemporalSchedule(h.Temporal, ctx.Context(), current.TemporalScheduleID); err != nil {
		ctx.L().Error("delete temporal schedule fail", zap.Error(err), zap.Uint64("id", req.ID))
		return res.FailMsg("删除 Temporal Schedule 失败")
	}
	jobSchedule := h.Q.JobSchedule
	info, err := jobSchedule.Where(jobSchedule.ID.Eq(req.ID), jobSchedule.Status.Neq(models.JobScheduleStatusDeleted)).
		Update(jobSchedule.Status, models.JobScheduleStatusDeleted)
	if err != nil {
		return res.FailDefault
	}
	if info.RowsAffected == 0 {
		return res.FailMsg("任务不存在")
	}
	return nil
}

// @Summary 启用或停用任务调度配置
// @Tags JobSchedule
// @Accept json
// @Produce json
// @Param req body ReqJobScheduleSwitch true "切换参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/sys/job/schedule/switch [post]
func (h *JobScheduleHandler) Switch(ctx *handler.Ctx, req *ReqJobScheduleSwitch) error {
	current, err := findActiveJobSchedule(h.Q, req.ID)
	if err != nil {
		return err
	}
	status := models.JobScheduleStatusDisabled
	if req.Enabled {
		status = models.JobScheduleStatusEnabled
	}
	if err = switchTemporalSchedule(h.Temporal, ctx.Context(), current.TemporalScheduleID, req.Enabled); err != nil {
		if isTemporalNotFound(err) {
			inputValue, inputErr := parseJSONValue(current.InputJSON)
			if inputErr != nil {
				return res.FailMsg("输入参数不是合法 JSON")
			}
			if syncErr := syncTemporalSchedule(h.Temporal, ctx.Context(), current, inputValue, false); syncErr != nil {
				ctx.L().Error("sync missing temporal schedule fail", zap.Error(syncErr), zap.Uint64("id", req.ID))
				return res.FailMsg("同步 Temporal Schedule 失败")
			}
			err = switchTemporalSchedule(h.Temporal, ctx.Context(), current.TemporalScheduleID, req.Enabled)
		}
	}
	if err != nil {
		ctx.L().Error("switch temporal schedule fail", zap.Error(err), zap.Uint64("id", req.ID))
		return res.FailMsg("切换 Temporal Schedule 失败")
	}
	jobSchedule := h.Q.JobSchedule
	info, err := jobSchedule.Where(jobSchedule.ID.Eq(req.ID), jobSchedule.Status.Neq(models.JobScheduleStatusDeleted)).
		Update(jobSchedule.Status, status)
	if err != nil {
		return res.FailDefault
	}
	if info.RowsAffected == 0 {
		return res.FailMsg("任务不存在")
	}
	return nil
}

// @Summary 同步任务调度配置到 Temporal
// @Tags JobSchedule
// @Accept json
// @Produce json
// @Param req body ReqJobScheduleID true "任务ID"
// @Success 200 {object} res.Response "成功"
// @Router /api/sys/job/schedule/sync [post]
func (h *JobScheduleHandler) Sync(ctx *handler.Ctx, req *ReqJobScheduleID) error {
	current, err := findActiveJobSchedule(h.Q, req.ID)
	if err != nil {
		return err
	}
	inputValue, err := parseJSONValue(current.InputJSON)
	if err != nil {
		return res.FailMsg("输入参数不是合法 JSON")
	}
	if err = syncTemporalSchedule(h.Temporal, ctx.Context(), current, inputValue, false); err != nil {
		ctx.L().Error("sync temporal schedule fail", zap.Error(err), zap.Uint64("id", req.ID))
		return res.FailMsg("同步 Temporal Schedule 失败")
	}
	return nil
}

// @Summary 立即触发任务调度
// @Tags JobSchedule
// @Accept json
// @Produce json
// @Param req body ReqJobScheduleID true "任务ID"
// @Success 200 {object} res.Response "成功"
// @Router /api/sys/job/schedule/trigger [post]
func (h *JobScheduleHandler) Trigger(ctx *handler.Ctx, req *ReqJobScheduleID) error {
	current, err := findActiveJobSchedule(h.Q, req.ID)
	if err != nil {
		return err
	}
	if err = triggerTemporalSchedule(h.Temporal, ctx.Context(), current.TemporalScheduleID); err != nil {
		ctx.L().Error("trigger temporal schedule fail", zap.Error(err), zap.Uint64("id", req.ID))
		return res.FailMsg("触发 Temporal Schedule 失败")
	}
	return nil
}

func (req *ReqJobScheduleCreate) normalize() {
	req.JobCode = strings.TrimSpace(req.JobCode)
	req.JobName = strings.TrimSpace(req.JobName)
	req.WorkflowType = strings.TrimSpace(req.WorkflowType)
	req.TaskQueue = strings.TrimSpace(req.TaskQueue)
	req.ScheduleType = strings.ToUpper(strings.TrimSpace(req.ScheduleType))
	req.CronExpr = strings.TrimSpace(req.CronExpr)
	req.InputJSON = strings.TrimSpace(req.InputJSON)
	req.Status = strings.ToUpper(strings.TrimSpace(req.Status))
	req.TemporalScheduleID = strings.TrimSpace(req.TemporalScheduleID)
	req.TemporalWorkflowIDPrefix = strings.TrimSpace(req.TemporalWorkflowIDPrefix)
	req.Description = strings.TrimSpace(req.Description)
}

func (req *ReqJobScheduleUpdate) normalize() {
	str.TrimStringPtr(req.JobName, nil)
	str.TrimStringPtr(req.WorkflowType, nil)
	str.TrimStringPtr(req.TaskQueue, nil)
	str.TrimStringPtr(req.ScheduleType, func(value string) string { return strings.ToUpper(value) })
	str.TrimStringPtr(req.CronExpr, nil)
	str.TrimStringPtr(req.InputJSON, nil)
	str.TrimStringPtr(req.Status, func(value string) string { return strings.ToUpper(value) })
	str.TrimStringPtr(req.TemporalWorkflowIDPrefix, nil)
	str.TrimStringPtr(req.Description, nil)
}

func (req *ReqJobScheduleCreate) toModel() (*models.JobSchedule, any, error) {
	inputJSON, inputValue, err := normalizeJSONText(req.InputJSON)
	if err != nil {
		return nil, nil, err
	}
	status := req.Status
	if status == "" {
		status = models.JobScheduleStatusEnabled
	}
	scheduleID := req.TemporalScheduleID
	if scheduleID == "" {
		scheduleID = req.JobCode
	}
	workflowIDPrefix := req.TemporalWorkflowIDPrefix
	if workflowIDPrefix == "" {
		workflowIDPrefix = req.JobCode
	}
	return &models.JobSchedule{
		JobCode:                  req.JobCode,
		JobName:                  req.JobName,
		WorkflowType:             req.WorkflowType,
		TaskQueue:                req.TaskQueue,
		ScheduleType:             req.ScheduleType,
		CronExpr:                 req.CronExpr,
		IntervalSeconds:          req.IntervalSeconds,
		StartTime:                req.StartTime,
		EndTime:                  req.EndTime,
		InputJSON:                inputJSON,
		Status:                   status,
		TemporalScheduleID:       scheduleID,
		TemporalWorkflowIDPrefix: workflowIDPrefix,
		Description:              req.Description,
	}, inputValue, nil
}

func (req *ReqJobScheduleUpdate) mergeCurrent(q *query.Query) (*models.JobSchedule, *models.JobSchedule, any, error) {
	current, err := findActiveJobSchedule(q, req.ID)
	if err != nil {
		return nil, nil, nil, err
	}
	next := *current
	if req.JobName != nil {
		next.JobName = *req.JobName
	}
	if req.WorkflowType != nil {
		next.WorkflowType = *req.WorkflowType
	}
	if req.TaskQueue != nil {
		next.TaskQueue = *req.TaskQueue
	}
	if req.ScheduleType != nil {
		next.ScheduleType = *req.ScheduleType
	}
	if req.CronExpr != nil {
		next.CronExpr = *req.CronExpr
	}
	if req.IntervalSeconds != nil {
		next.IntervalSeconds = req.IntervalSeconds
	}
	if req.StartTime != nil {
		next.StartTime = req.StartTime
	}
	if req.EndTime != nil {
		next.EndTime = req.EndTime
	}
	if req.Status != nil {
		next.Status = *req.Status
	}
	if req.TemporalWorkflowIDPrefix != nil {
		next.TemporalWorkflowIDPrefix = *req.TemporalWorkflowIDPrefix
	}
	if req.Description != nil {
		next.Description = *req.Description
	}
	if next.TemporalScheduleID == "" {
		next.TemporalScheduleID = next.JobCode
	}
	if next.TemporalWorkflowIDPrefix == "" {
		next.TemporalWorkflowIDPrefix = next.JobCode
	}

	inputValue, err := parseJSONValue(next.InputJSON)
	if req.InputJSON != nil {
		var inputJSON datatypes.JSON
		inputJSON, inputValue, err = normalizeJSONText(*req.InputJSON)
		next.InputJSON = inputJSON
	}
	if err != nil {
		return nil, nil, nil, err
	}
	return current, &next, inputValue, nil
}

func validateJobSchedule(m *models.JobSchedule) error {
	if m.JobCode == "" {
		return res.FailMsg("任务编码不能为空")
	}
	if m.JobName == "" {
		return res.FailMsg("任务名称不能为空")
	}
	if m.WorkflowType == "" {
		return res.FailMsg("Workflow 类型不能为空")
	}
	if m.TaskQueue == "" {
		return res.FailMsg("Task Queue 不能为空")
	}
	if _, ok := validJobScheduleTypes[m.ScheduleType]; !ok {
		return res.FailMsg("调度类型无效")
	}
	if _, ok := validJobScheduleStatuses[m.Status]; !ok {
		return res.FailMsg("任务状态无效")
	}
	if m.ScheduleType == models.JobScheduleTypeCron && strings.TrimSpace(m.CronExpr) == "" {
		return res.FailMsg("CRON 调度必须填写 cron 表达式")
	}
	if m.ScheduleType == models.JobScheduleTypeInterval && (m.IntervalSeconds == nil || *m.IntervalSeconds <= 0) {
		return res.FailMsg("INTERVAL 调度必须填写大于 0 的间隔秒数")
	}
	if m.ScheduleType == models.JobScheduleTypeOnce && m.StartTime == nil {
		return res.FailMsg("ONCE 调度必须填写开始时间")
	}
	if m.StartTime != nil && m.EndTime != nil && !m.EndTime.After(*m.StartTime) {
		return res.FailMsg("结束时间必须晚于开始时间")
	}
	return nil
}

var validJobScheduleTypes = map[string]struct{}{
	models.JobScheduleTypeOnce:     {},
	models.JobScheduleTypeCron:     {},
	models.JobScheduleTypeInterval: {},
}

var validJobScheduleStatuses = map[string]struct{}{
	models.JobScheduleStatusEnabled:  {},
	models.JobScheduleStatusDisabled: {},
	models.JobScheduleStatusDeleted:  {},
}

func buildJobScheduleOptions(items map[string]string) []RespJobScheduleOption {
	keys := make([]string, 0, len(items))
	for value := range items {
		if strings.TrimSpace(value) == "" {
			continue
		}
		keys = append(keys, value)
	}
	sort.Strings(keys)

	options := make([]RespJobScheduleOption, 0, len(keys))
	for _, value := range keys {
		label := strings.TrimSpace(items[value])
		if label == "" {
			label = value
		}
		options = append(options, RespJobScheduleOption{Label: label, Value: value})
	}
	return options
}

func findActiveJobSchedule(q *query.Query, id uint64) (*models.JobSchedule, error) {
	jobSchedule := q.JobSchedule
	current, err := jobSchedule.Where(jobSchedule.ID.Eq(id), jobSchedule.Status.Neq(models.JobScheduleStatusDeleted)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, res.FailMsg("任务不存在")
		}
		return nil, res.FailDefault
	}
	return current, nil
}

func syncTemporalSchedule(temporal service.TemporalService, ctx context.Context, m *models.JobSchedule, inputValue any, createOnly bool) error {
	
	options, err := buildScheduleOptions(m, inputValue)
	if err != nil {
		return err
	}
	if createOnly {
		_, err = temporal.CreateSchedule(ctx, options)
		return err
	}
	err = temporal.UpdateSchedule(ctx, options.ID, client.ScheduleUpdateOptions{
		DoUpdate: func(client.ScheduleUpdateInput) (*client.ScheduleUpdate, error) {
			spec := options.Spec
			return &client.ScheduleUpdate{
				Schedule: &client.Schedule{
					Action: options.Action,
					Spec:   &spec,
					Policy: &client.SchedulePolicies{},
					State: &client.ScheduleState{
						Paused:           options.Paused,
						LimitedActions:   options.RemainingActions > 0,
						RemainingActions: options.RemainingActions,
					},
				},
			}, nil
		},
	})
	if err == nil {
		return nil
	}
	_, err = temporal.CreateSchedule(ctx, options)
	return err
}

func deleteTemporalSchedule(temporal service.TemporalService, ctx context.Context, scheduleID string) error {
	if scheduleID == "" {
		return nil
	}
	
	return temporal.DeleteSchedule(ctx, scheduleID)
}

func switchTemporalSchedule(temporal service.TemporalService, ctx context.Context, scheduleID string, enabled bool) error {
	
	if enabled {
		return temporal.UnpauseSchedule(ctx, scheduleID, client.ScheduleUnpauseOptions{})
	}
	return temporal.PauseSchedule(ctx, scheduleID, client.SchedulePauseOptions{Note: "disabled by admin"})
}

func triggerTemporalSchedule(temporal service.TemporalService, ctx context.Context, scheduleID string) error {
	
	return temporal.TriggerSchedule(ctx, scheduleID, client.ScheduleTriggerOptions{})
}

func isTemporalNotFound(err error) bool {
	var notFound *serviceerror.NotFound
	return errors.As(err, &notFound)
}

func buildScheduleOptions(m *models.JobSchedule, inputValue any) (client.ScheduleOptions, error) {
	spec, remainingActions, err := buildScheduleSpec(m)
	if err != nil {
		return client.ScheduleOptions{}, err
	}
	return client.ScheduleOptions{
		ID:   m.TemporalScheduleID,
		Spec: spec,
		Action: &client.ScheduleWorkflowAction{
			ID:        fmt.Sprintf("%s-dispatch", m.TemporalWorkflowIDPrefix),
			Workflow:  temporaljob.DispatchWorkflowName,
			TaskQueue: m.TaskQueue,
			Args: []any{temporaljob.DispatchInput{
				JobCode:          m.JobCode,
				WorkflowType:     m.WorkflowType,
				TaskQueue:        m.TaskQueue,
				WorkflowIDPrefix: m.TemporalWorkflowIDPrefix,
				Input:            inputValue,
			}},
		},
		Paused:           m.Status != models.JobScheduleStatusEnabled,
		RemainingActions: remainingActions,
	}, nil
}

func buildScheduleSpec(m *models.JobSchedule) (client.ScheduleSpec, int, error) {
	spec := client.ScheduleSpec{}
	if m.StartTime != nil {
		spec.StartAt = *m.StartTime
	}
	if m.EndTime != nil {
		spec.EndAt = *m.EndTime
	}
	switch m.ScheduleType {
	case models.JobScheduleTypeCron:
		spec.CronExpressions = []string{m.CronExpr}
	case models.JobScheduleTypeInterval:
		spec.Intervals = []client.ScheduleIntervalSpec{{Every: time.Duration(*m.IntervalSeconds) * time.Second}}
	case models.JobScheduleTypeOnce:
		start := m.StartTime.UTC()
		spec.TimeZoneName = "UTC"
		spec.Calendars = []client.ScheduleCalendarSpec{{
			Second:     []client.ScheduleRange{{Start: start.Second()}},
			Minute:     []client.ScheduleRange{{Start: start.Minute()}},
			Hour:       []client.ScheduleRange{{Start: start.Hour()}},
			DayOfMonth: []client.ScheduleRange{{Start: start.Day()}},
			Month:      []client.ScheduleRange{{Start: int(start.Month())}},
			Year:       []client.ScheduleRange{{Start: start.Year()}},
		}}
		return spec, 1, nil
	default:
		return client.ScheduleSpec{}, 0, res.FailMsg("调度类型无效")
	}
	return spec, 0, nil
}


func normalizeJSONText(text string) (datatypes.JSON, any, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, nil, nil
	}
	var value any
	if err := sonic.UnmarshalString(text, &value); err != nil {
		return nil, nil, res.FailMsg("输入参数不是合法 JSON")
	}
	data, err := sonic.Marshal(value)
	if err != nil {
		return nil, nil, res.FailDefault
	}
	return datatypes.JSON(data), value, nil
}

func parseJSONValue(data datatypes.JSON) (any, error) {
	if len(data) == 0 {
		return nil, nil
	}
	var value any
	if err := sonic.Unmarshal(data, &value); err != nil {
		return nil, err
	}
	return value, nil
}

func appendPagingQuery(req *v1.PagingRequest, extra map[string]any) {
	current := strings.TrimSpace(req.GetQuery())
	if current == "" {
		data, _ := sonic.MarshalString(extra)
		req.FilteringType = &v1.PagingRequest_Query{Query: data}
		return
	}
	req.FilteringType = &v1.PagingRequest_Query{Query: fmt.Sprintf(`{"$and":[%s,%s]}`, current, mustJSON(extra))}
}

func mustJSON(value any) string {
	data, _ := sonic.MarshalString(value)
	return data
}
