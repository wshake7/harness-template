package models

import (
	"gorm.io/plugin/soft_delete"
	"orm-crud/gormc/mixin"
)

func init() {
	Models = append(Models, &SysResourceApi{})
}

// SysResourceApi maps to sys_resource_api.
type SysResourceApi struct {
	mixin.AutoIncrementID
	mixin.CreatedAt
	mixin.UpdatedAt
	mixin.SortOrder
	mixin.OperatorID
	mixin.Remark
	mixin.IsEnabled
	DeletedAt soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;not null;default:0;index" json:"deletedAt"`
	Module    string                `gorm:"column:module;type:varchar(128);not null;default:'';comment:所属业务模块" json:"module"`
	Path      string                `gorm:"column:path;type:varchar(255);not null;uniqueIndex:idx_sys_resource_api_method_path_active,priority:2,where:deleted_at = 0;comment:接口路径模板" json:"path"`
	Method    string                `gorm:"column:method;type:varchar(16);not null;uniqueIndex:idx_sys_resource_api_method_path_active,priority:1,where:deleted_at = 0;comment:请求方法" json:"method"`
}

func (SysResourceApi) TableName() string {
	return "sys_resource_api"
}
