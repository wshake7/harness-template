package models

import (
	"orm-crud/gormc/mixin"

	"gorm.io/plugin/soft_delete"
)

func init() {
	Models = append(Models, &SysDictEntry{})
}

type SysDictEntry struct {
	mixin.AutoIncrementID
	mixin.CreatedAt
	mixin.UpdatedAt
	mixin.OperatorID
	mixin.SortOrder
	mixin.IsEnabled
	mixin.Remark
	DeletedAt      soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;not null;default:0;index:idx_sys_dict_entry_deleted_at" json:"deletedAt"`
	LabelComponent string                `gorm:"column:label_component;type:varchar(255);not null;default:'';comment:字典项的显示标签组件" json:"labelComponent"`
	EntryLabel     string                `gorm:"column:entry_label;type:varchar(255);not null;comment:字典项的显示标签" json:"entryLabel"`
	EntryValue     string                `gorm:"column:entry_value;type:varchar(255);not null;uniqueIndex:idx_sys_dict_entry_type_lang_value_active,priority:3,where:deleted_at = 0;comment:字典项的实际值" json:"entryValue"`
	LanguageCode   string                `gorm:"column:language_code;type:varchar(32);not null;default:'';uniqueIndex:idx_sys_dict_entry_type_lang_value_active,priority:2,where:deleted_at = 0;comment:语言代码" json:"languageCode"`
	SysDictTypeId  uint64                `gorm:"column:sys_dict_type_id;type:bigint;not null;index:idx_sys_dict_entry_dict_type_id;uniqueIndex:idx_sys_dict_entry_type_lang_value_active,priority:1,where:deleted_at = 0;comment:字典类型ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"sysDictTypeId"`
	SysDictType    *SysDictType          `gorm:"foreignKey:SysDictTypeId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"sysDictType"`
}

func (SysDictEntry) TableName() string {
	return "sys_dict_entry"
}
