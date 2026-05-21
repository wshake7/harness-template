package models

import (
	"time"

	"orm-crud/gormc/mixin"

	"gorm.io/datatypes"
)

const (
	JobScheduleTypeOnce       = "ONCE"
	JobScheduleTypeCron       = "CRON"
	JobScheduleTypeInterval   = "INTERVAL"
	JobScheduleStatusEnabled  = "ENABLED"
	JobScheduleStatusDisabled = "DISABLED"
	JobScheduleStatusDeleted  = "DELETED"
)

func init() {
	Models = append(Models, &JobSchedule{})
}

type JobSchedule struct {
	mixin.AutoIncrementID
	mixin.CreatedAt
	mixin.UpdatedAt
	JobCode                  string         `gorm:"column:job_code;type:varchar(128);not null;uniqueIndex;comment:业务任务编码" json:"jobCode"`
	JobName                  string         `gorm:"column:job_name;type:varchar(255);not null;comment:任务名称" json:"jobName"`
	WorkflowType             string         `gorm:"column:workflow_type;type:varchar(255);not null;comment:Temporal Workflow 类型" json:"workflowType"`
	TaskQueue                string         `gorm:"column:task_queue;type:varchar(255);not null;comment:Temporal Task Queue" json:"taskQueue"`
	ScheduleType             string         `gorm:"column:schedule_type;type:varchar(32);not null;index:idx_schedule_type;comment:ONCE/CRON/INTERVAL" json:"scheduleType"`
	CronExpr                 string         `gorm:"column:cron_expr;type:varchar(128);default:'';comment:cron 表达式" json:"cronExpr"`
	IntervalSeconds          *int           `gorm:"column:interval_seconds;type:int;default:null;comment:间隔秒数" json:"intervalSeconds"`
	StartTime                *time.Time     `gorm:"column:start_time;default:null;comment:开始时间" json:"startTime"`
	EndTime                  *time.Time     `gorm:"column:end_time;default:null;comment:结束时间" json:"endTime"`
	InputJSON                datatypes.JSON `gorm:"column:input_json;type:json;default:null;comment:Workflow 输入参数" json:"inputJSON"`
	Status                   string         `gorm:"column:status;type:varchar(32);not null;default:'ENABLED';index:idx_job_schedule_status;comment:ENABLED/DISABLED/DELETED" json:"status"`
	TemporalScheduleID       string         `gorm:"column:temporal_schedule_id;type:varchar(255);default:'';comment:Temporal Schedule ID" json:"temporalScheduleID"`
	TemporalWorkflowIDPrefix string         `gorm:"column:temporal_workflow_id_prefix;type:varchar(255);default:'';comment:Workflow ID 前缀" json:"temporalWorkflowIDPrefix"`
	Description              string         `gorm:"column:description;type:varchar(512);default:''" json:"description"`
}

func (*JobSchedule) TableName() string {
	return "job_schedule"
}
