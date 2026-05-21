package models

import (
	"orm-crud/gormc/mixin"

	"gorm.io/plugin/soft_delete"
)

func init() {
	Models = append(Models, &SysRoleApi{})
}

type SysRoleApi struct {
	mixin.AutoIncrementID
	mixin.CreatedAt
	mixin.UpdatedAt
	mixin.OperatorID
	RoleID         uint64                `gorm:"column:role_id;type:bigint;not null;index:idx_sys_role_api_role_id;uniqueIndex:idx_sys_role_api_role_api_active,priority:1,where:deleted_at = 0;comment:角色ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"roleID"`
	ApiID          uint64                `gorm:"column:api_id;type:bigint;not null;index:idx_sys_role_api_api_id;uniqueIndex:idx_sys_role_api_role_api_active,priority:2,where:deleted_at = 0;comment:API资源ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"apiID"`
	DeletedAt      soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;not null;default:0;index:idx_sys_role_api_deleted_at" json:"deletedAt"`
	SysRole        *SysRole              `gorm:"foreignKey:RoleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"sysRole"`
	SysResourceApi *SysResourceApi       `gorm:"foreignKey:ApiID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"sysResourceApi"`
}

func (SysRoleApi) TableName() string {
	return "sys_role_api"
}
