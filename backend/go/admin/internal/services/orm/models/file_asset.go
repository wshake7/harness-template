package models

import (
	"orm-crud/gormc/mixin"

	"gorm.io/datatypes"
	"gorm.io/plugin/soft_delete"
)

const (
	FileAssetStatusActive  = "active"
	FileAssetStatusDeleted = "deleted"
)

func init() {
	Models = append(Models, &FileAsset{})
}

type FileAsset struct {
	mixin.AutoIncrementID
	mixin.CreatedAt
	mixin.UpdatedAt
	mixin.OperatorID
	mixin.Remark
	DeletedAt    soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;not null;default:0;index:idx_file_asset_deleted_at" json:"deletedAt"`
	Engine       string                `gorm:"column:engine;type:varchar(32);not null;comment:存储引擎" json:"engine"`
	Bucket       string                `gorm:"column:bucket;type:varchar(128);not null;comment:bucket" json:"bucket"`
	ObjectKey    string                `gorm:"column:object_key;type:varchar(512);not null;uniqueIndex:uk_file_asset_object_key_active,where:deleted_at = 0;comment:对象key" json:"objectKey"`
	OriginalName string                `gorm:"column:original_name;type:varchar(255);not null;comment:原始文件名" json:"originalName"`
	ContentType  string                `gorm:"column:content_type;type:varchar(128);not null;default:'';comment:内容类型" json:"contentType"`
	Extension    string                `gorm:"column:extension;type:varchar(32);not null;default:'';comment:文件扩展名" json:"extension"`
	Size         int64                 `gorm:"column:size;type:bigint;not null;default:0;comment:文件大小" json:"size"`
	SHA256       string                `gorm:"column:sha256;type:char(64);not null;default:'';index:idx_file_asset_sha256;comment:SHA256" json:"sha256"`
	BizType      string                `gorm:"column:biz_type;type:varchar(64);not null;default:'';index:idx_file_asset_biz,priority:1;comment:业务类型" json:"bizType"`
	BizID        string                `gorm:"column:biz_id;type:varchar(128);not null;default:'';index:idx_file_asset_biz,priority:2;comment:业务ID" json:"bizID"`
	Metadata     datatypes.JSONMap     `gorm:"column:metadata;default:'{}';comment:扩展元数据" json:"metadata"`
	Status       string                `gorm:"column:status;type:varchar(32);not null;default:'active';comment:状态" json:"status"`
}

func (*FileAsset) TableName() string {
	return "file_asset"
}
