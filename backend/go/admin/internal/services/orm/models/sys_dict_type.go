package models

import (
	"orm-crud/gormc/mixin"

	"gorm.io/plugin/soft_delete"
)

func init() {
	Models = append(Models, &SysDictType{})
}

type SysDictType struct {
	mixin.AutoIncrementID
	mixin.CreatedAt
	mixin.UpdatedAt
	mixin.OperatorID
	mixin.IsEnabled
	mixin.SortOrder
	mixin.Remark
	DeletedAt soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;not null;default:0;index:idx_sys_dict_type_deleted_at" json:"deletedAt"`
	TypeCode  string                `gorm:"column:type_code;type:varchar(128);not null;uniqueIndex:idx_sys_dict_type_type_code_active,where:deleted_at = 0;comment:字典类型唯一代码" json:"typeCode"`
	TypeName  string                `gorm:"column:type_name;type:varchar(255);not null;comment:字典类型名称" json:"typeName"`
	Entries   []SysDictEntry        `gorm:"foreignKey:SysDictTypeId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"entries"`
}

func (SysDictType) TableName() string {
	return "sys_dict_type"
}
