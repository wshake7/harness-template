package models

import (
	"orm-crud/gormc/mixin"

	"gorm.io/plugin/soft_delete"
)

func init() {
	Models = append(Models, &SysLanguageEntry{})
}

type SysLanguageEntry struct {
	mixin.AutoIncrementID
	mixin.CreatedAt
	mixin.UpdatedAt
	mixin.OperatorID
	mixin.SortOrder
	mixin.IsEnabled
	mixin.Remark
	DeletedAt         soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;not null;default:0;index:idx_sys_language_entry_deleted_at" json:"deletedAt"`
	EntryCode         string                `gorm:"column:entry_code;type:varchar(128);not null;uniqueIndex:idx_sys_language_entry_code_type_active,priority:1,where:deleted_at = 0;comment:语言条目编码" json:"entryCode"`
	EntryValue        string                `gorm:"column:entry_value;type:varchar(255);not null;comment:语言值" json:"entryValue"`
	SysLanguageTypeId uint64                `gorm:"column:sys_language_type_id;type:bigint;not null;uniqueIndex:idx_sys_language_entry_code_type_active,priority:2,where:deleted_at = 0;comment:语言类型ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"sysLanguageTypeId"`
	SysLanguageType   *SysLanguageType      `gorm:"foreignKey:SysLanguageTypeId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"sysLanguageType"`
}

func (SysLanguageEntry) TableName() string {
	return "sys_language_entry"
}
