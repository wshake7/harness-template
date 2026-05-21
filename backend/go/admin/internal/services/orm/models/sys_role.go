package models

import (
	"orm-crud/gormc/mixin"

	"gorm.io/plugin/soft_delete"
)

func init() {
	Models = append(Models, &SysRole{})
}

type SysRole struct {
	mixin.AutoIncrementID
	mixin.CreatedAt
	mixin.UpdatedAt
	mixin.Remark
	mixin.OperatorID
	mixin.IsEnabled
	DeletedAt     soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;not null;default:0;index:idx_sys_role_deleted_at" json:"deletedAt"`
	Name          string                `gorm:"column:name;type:varchar(255);not null;comment:角色名称" json:"name"`
	Code          string                `gorm:"column:code;type:varchar(128);not null;uniqueIndex:idx_sys_role_code_active,where:deleted_at = 0;comment:角色标识" json:"code"`
	ParentID      *uint64               `gorm:"column:parent_id;type:bigint;comment:父级ID" json:"parentID"`
	ParentSysRole *SysRole              `gorm:"foreignKey:ParentID;references:ID" json:"parentSysRole"`
	Children      []SysRole             `gorm:"foreignKey:ParentID;references:ID" json:"children"`
}

func (*SysRole) TableName() string {
	return "sys_role"
}
