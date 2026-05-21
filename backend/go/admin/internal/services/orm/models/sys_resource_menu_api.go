package models

import (
	"orm-crud/gormc/mixin"

	"gorm.io/plugin/soft_delete"
)

func init() {
	Models = append(Models, &SysResourceMenuApi{})
}

type SysResourceMenuApi struct {
	mixin.AutoIncrementID
	mixin.CreatedAt
	mixin.UpdatedAt
	mixin.OperatorID
	MenuID          uint64                `gorm:"column:menu_id;type:bigint;not null;index:idx_sys_resource_menu_api_menu_id;uniqueIndex:idx_sys_resource_menu_api_menu_api_active,priority:1,where:deleted_at = 0;comment:菜单资源ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"menuID"`
	ApiID           uint64                `gorm:"column:api_id;type:bigint;not null;index:idx_sys_resource_menu_api_api_id;uniqueIndex:idx_sys_resource_menu_api_menu_api_active,priority:2,where:deleted_at = 0;comment:API资源ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"apiID"`
	DeletedAt       soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;not null;default:0;index:idx_sys_resource_menu_api_deleted_at" json:"deletedAt"`
	SysResourceMenu *SysResourceMenu      `gorm:"foreignKey:MenuID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"sysResourceMenu"`
	SysResourceApi  *SysResourceApi       `gorm:"foreignKey:ApiID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"sysResourceApi"`
}

func (SysResourceMenuApi) TableName() string {
	return "sys_resource_menu_api"
}
