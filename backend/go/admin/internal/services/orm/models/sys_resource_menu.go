package models

import (
	"gorm.io/plugin/soft_delete"
	"orm-crud/gormc/mixin"
)

func init() {
	Models = append(Models, &SysResourceMenu{})
}

// SysResourceMenu 对应表 sys_resource_menu
type SysResourceMenu struct {
	mixin.AutoIncrementID
	mixin.CreatedAt
	mixin.UpdatedAt
	mixin.OperatorID
	mixin.Remark
	mixin.SortOrder
	mixin.Metadata
	mixin.IsEnabled
	DeletedAt soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;not null;default:0;index" json:"deletedAt"`
	MenuType  string                `gorm:"column:menu_type;type:varchar(32);not null;default:'MENU';comment:菜单类型 CATALOG: 目录 MENU: 菜单 BUTTON: 按钮 EMBEDDED: 内嵌 LINK: 外链" json:"menuType"`
	Path      string                `gorm:"column:path;type:varchar(1024);not null;comment:路径，类型为按钮时为操作名" json:"path"`
	Redirect  string                `gorm:"column:redirect;type:varchar(1024);not null;default:'';comment:重定向地址" json:"redirect"`
	Alias     string                `gorm:"column:alias;type:varchar(255);not null;default:'';comment:路由别名" json:"alias"`
	Name      string                `gorm:"column:name;type:varchar(255);not null;default:'';comment:路由命名" json:"name"`
	Component string                `gorm:"column:component;type:varchar(255);not null;default:'';comment:前端页面组件" json:"component"`

	// 简单树结构字段（保留父级关系与路径）
	ParentID *uint64 `gorm:"column:parent_id;type:bigint;comment:父级ID" json:"parentID"`
	TreePath *string `gorm:"column:tree_path;type:varchar(1024);comment:节点路径" json:"treePath"`

	ParentSysResourceMenu *SysResourceMenu `gorm:"foreignKey:ParentID;references:ID" json:"parentSysResourceMenu"`
}

func (SysResourceMenu) TableName() string {
	return "sys_resource_menu"
}

type CatalogMeta struct {
	Icon   string `json:"icon"`
	Hidden bool   `json:"hidden"`
}

type MenuMeta struct {
	Icon   string `json:"icon"`
	Hidden bool   `json:"hidden"`
}

type ButtonMeta struct {
	Authorities []string `json:"authorities"`
}

type EMBEDDEDMeta struct {
	Authorities []string `json:"authorities"`
}

type LinkMeta struct {
	Authorities []string `json:"authorities"`
}
