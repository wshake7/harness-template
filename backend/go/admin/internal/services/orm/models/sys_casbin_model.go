package models

import (
	"orm-crud/gormc/mixin"

	"gorm.io/plugin/soft_delete"
)

func init() {
	Models = append(Models, &SysCasbinModel{})
}

type SysCasbinModel struct {
	mixin.AutoIncrementID
	mixin.CreatedAt
	mixin.UpdatedAt
	mixin.OperatorID
	mixin.IsEnabled
	mixin.Remark
	DeletedAt soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;not null;default:0;index:idx_sys_casbin_model_deleted_at" json:"deletedAt"`
	Name      string                `gorm:"column:name;type:varchar(255);not null;uniqueIndex:idx_sys_casbin_model_name_active,where:deleted_at = 0;comment:模型名称" json:"name"`
	Content   string                `gorm:"column:content;type:text;not null;comment:模型内容" json:"content"`
}

func (SysCasbinModel) TableName() string {
	return "sys_casbin_model"
}
