package models

import (
	"orm-crud/gormc/mixin"

	"gorm.io/plugin/soft_delete"
)

func init() {
	Models = append(Models, &SysLanguageType{})
}

type SysLanguageType struct {
	mixin.AutoIncrementID
	mixin.CreatedAt
	mixin.UpdatedAt
	mixin.OperatorID
	mixin.SortOrder
	mixin.IsEnabled
	DeletedAt soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;not null;default:0;index:idx_sys_language_type_deleted_at" json:"deletedAt"`
	TypeCode  string                `gorm:"column:type_code;type:varchar(128);not null;uniqueIndex:idx_sys_language_type_code_active,where:deleted_at = 0;comment:标准语言代码" json:"typeCode"`
	TypeName  string                `gorm:"column:type_name;type:varchar(255);not null;comment:语言名称" json:"typeName"`
	IsDefault bool                  `gorm:"column:is_default;default:false;uniqueIndex:idx_sys_language_type_default_active,where:is_default = true AND deleted_at = 0;comment:是否为默认语言" json:"isDefault"`
	Entries   []SysLanguageEntry    `gorm:"foreignKey:SysLanguageTypeId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"entries"`
}

func (SysLanguageType) TableName() string {
	return "sys_language_type"
}
