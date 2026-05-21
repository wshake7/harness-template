package models

import (
	"time"

	"orm-crud/gormc/mixin"

	"gorm.io/datatypes"
)

const (
	JobExecutionStatusRunning  = "RUNNING"
	JobExecutionStatusSuccess  = "SUCCESS"
	JobExecutionStatusFailed   = "FAILED"
	JobExecutionStatusCanceled = "CANCELED"
	JobExecutionStatusTimeout  = "TIMEOUT"
)

func init() {
	Models = append(Models, &JobExecution{})
}

type JobExecution struct {
	mixin.AutoIncrementID
	mixin.CreatedAt
	mixin.UpdatedAt
	JobCode            string         `gorm:"column:job_code;type:varchar(128);not null;index:idx_job_code" json:"jobCode"`
	TemporalWorkflowID string         `gorm:"column:temporal_workflow_id;type:varchar(255);not null;uniqueIndex:uk_workflow_run,priority:1" json:"temporalWorkflowID"`
	TemporalRunID      string         `gorm:"column:temporal_run_id;type:varchar(255);default:'';uniqueIndex:uk_workflow_run,priority:2" json:"temporalRunID"`
	TriggerTime        time.Time      `gorm:"column:trigger_time;not null;index:idx_trigger_time;comment:计划触发时间" json:"triggerTime"`
	StartTime          *time.Time     `gorm:"column:start_time;default:null" json:"startTime"`
	EndTime            *time.Time     `gorm:"column:end_time;default:null" json:"endTime"`
	Status             string         `gorm:"column:status;type:varchar(32);not null;default:'RUNNING';index:idx_job_execution_status;comment:RUNNING/SUCCESS/FAILED/CANCELED/TIMEOUT" json:"status"`
	InputJSON          datatypes.JSON `gorm:"column:input_json;type:json;default:null" json:"inputJSON"`
	ResultJSON         datatypes.JSON `gorm:"column:result_json;type:json;default:null" json:"resultJSON"`
	ErrorMessage       string         `gorm:"column:error_message;type:text;default:''" json:"errorMessage"`
	RetryCount         int            `gorm:"column:retry_count;type:int;not null;default:0" json:"retryCount"`
}

func (*JobExecution) TableName() string {
	return "job_execution"
}
