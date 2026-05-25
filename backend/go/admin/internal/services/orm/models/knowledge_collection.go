package models

import (
	"orm-crud/gormc/mixin"

	"gorm.io/plugin/soft_delete"
)

func init() {
	Models = append(Models, &KnowledgeCollection{})
}

// KnowledgeCollection 知识库集合管理
// 对应表 knowledge_collection，用于管理不同的知识库集合（Collection）
type KnowledgeCollection struct {
	mixin.AutoIncrementID
	mixin.CreatedAt
	mixin.UpdatedAt
	mixin.OperatorID
	mixin.IsEnabled
	mixin.Remark
	DeletedAt      soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;not null;default:0;index:idx_knowledge_collection_deleted_at" json:"deletedAt"`
	CollectionName string                `gorm:"column:collection_name;type:varchar(128);not null;uniqueIndex:idx_knowledge_collection_name_active,where:deleted_at = 0;comment:集合名称（对应 Milvus Collection）" json:"collectionName"`
	DisplayName    string                `gorm:"column:display_name;type:varchar(255);not null;comment:显示名称" json:"displayName"`
	MetricType     string                `gorm:"column:metric_type;type:varchar(32);not null;default:'COSINE';comment:相似度度量类型" json:"metricType"`
	IndexType      string                `gorm:"column:index_type;type:varchar(32);not null;default:'auto';comment:索引类型" json:"indexType"`
}

func (*KnowledgeCollection) TableName() string {
	return "knowledge_collection"
}
