package models

import (
	"orm-crud/gormc/mixin"

	"gorm.io/plugin/soft_delete"
)

func init() {
	Models = append(Models, &SysRoleMenu{})
}

type SysRoleMenu struct {
	mixin.AutoIncrementID
	mixin.CreatedAt
	mixin.UpdatedAt
	mixin.OperatorID
	RoleID          uint64                `gorm:"column:role_id;type:bigint;not null;index:idx_sys_role_menu_role_id;uniqueIndex:idx_sys_role_menu_role_menu_active,priority:1,where:deleted_at = 0;comment:角色ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"roleID"`
	MenuID          uint64                `gorm:"column:menu_id;type:bigint;not null;index:idx_sys_role_menu_menu_id;uniqueIndex:idx_sys_role_menu_role_menu_active,priority:2,where:deleted_at = 0;comment:菜单资源ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"menuID"`
	DeletedAt       soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;not null;default:0;index:idx_sys_role_menu_deleted_at" json:"deletedAt"`
	SysRole         *SysRole              `gorm:"foreignKey:RoleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"sysRole"`
	SysResourceMenu *SysResourceMenu      `gorm:"foreignKey:MenuID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"sysResourceMenu"`
}

func (SysRoleMenu) TableName() string {
	return "sys_role_menu"
}
