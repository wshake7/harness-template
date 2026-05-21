package models

import (
	"orm-crud/gormc/mixin"

	"gorm.io/plugin/soft_delete"
)

func init() {
	Models = append(Models, &SysUserRole{})
}

type SysUserRole struct {
	mixin.AutoIncrementID
	mixin.CreatedAt
	mixin.UpdatedAt
	mixin.OperatorID
	UserID    uint64                `gorm:"column:user_id;type:bigint;not null;index:idx_sys_user_role_user_id;uniqueIndex:idx_sys_user_role_user_role_active,priority:1,where:deleted_at = 0;comment:用户ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"userID"`
	RoleID    uint64                `gorm:"column:role_id;type:bigint;not null;index:idx_sys_user_role_role_id;uniqueIndex:idx_sys_user_role_user_role_active,priority:2,where:deleted_at = 0;comment:角色ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"roleID"`
	DeletedAt soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;not null;default:0;index:idx_sys_user_role_deleted_at" json:"deletedAt"`
	SysUser   *SysUser              `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"sysUser"`
	SysRole   *SysRole              `gorm:"foreignKey:RoleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"sysRole"`
}

func (SysUserRole) TableName() string {
	return "sys_user_role"
}
