package models

import (
	"orm-crud/gormc/mixin"

	"gorm.io/plugin/soft_delete"
)

func init() {
	Models = append(Models, &SysUser{})
}

type SysUser struct {
	mixin.AutoIncrementID
	mixin.CreatedAt
	mixin.UpdatedAt
	mixin.Remark
	mixin.OperatorID
	mixin.IsEnabled
	DeletedAt    soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;not null;default:0;index:idx_sys_user_deleted_at" json:"deletedAt"`
	Username     string                `gorm:"column:username;type:varchar(64);not null;uniqueIndex:idx_sys_user_username_active,where:deleted_at = 0;comment:用户名" json:"username"`
	Nickname     string                `gorm:"column:nickname;type:varchar(64);not null;default:'';comment:昵称" json:"nickname"`
	Password     string                `gorm:"column:password;type:varchar(255);not null;default:'';comment:密码" json:"-"`
	LanguageCode string                `gorm:"column:language_code;type:varchar(32);not null;default:'';comment:语言代码" json:"languageCode"`
}

func (*SysUser) TableName() string {
	return "sys_user"
}
