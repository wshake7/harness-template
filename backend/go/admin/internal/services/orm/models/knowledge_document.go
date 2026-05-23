package models

import (
	"orm-crud/gormc/mixin"

	"gorm.io/datatypes"
	"gorm.io/plugin/soft_delete"
)

func init() {
	Models = append(Models, &KnowledgeDocument{})
}

// KnowledgeDocument 知识库文档管理
// 对应表 knowledge_document，用于管理存入知识库的具体文档数据
type KnowledgeDocument struct {
	mixin.AutoIncrementID
	mixin.CreatedAt
	mixin.UpdatedAt
	mixin.OperatorID
	mixin.IsEnabled
	mixin.Remark
	DeletedAt          soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;not null;default:0;index:idx_knowledge_document_deleted_at" json:"deletedAt"`
	CollectionID       uint64                `gorm:"column:collection_id;type:bigint;not null;index:idx_knowledge_document_collection_id;comment:所属集合ID" json:"collectionID"`
	DocumentID         string                `gorm:"column:document_id;type:varchar(255);not null;uniqueIndex:idx_knowledge_document_doc_id_active,where:deleted_at = 0;comment:文档唯一标识（对应 Milvus ID）" json:"documentID"`
	Title              string                `gorm:"column:title;type:varchar(512);not null;comment:文档标题" json:"title"`
	Content            string                `gorm:"column:content;type:text;not null;comment:文档内容" json:"content"`
	ContentType        string                `gorm:"column:content_type;type:varchar(64);not null;default:'text';comment:内容类型 text/markdown/pdf/html" json:"contentType"`
	Source             string                `gorm:"column:source;type:varchar(512);default:'';comment:来源路径或URL" json:"source"`
	ChunkIndex         int                   `gorm:"column:chunk_index;type:int;not null;default:0;comment:分块索引，同一文档的多块内容" json:"chunkIndex"`
	TotalChunks        int                   `gorm:"column:total_chunks;type:int;not null;default:1;comment:总分块数" json:"totalChunks"`
	VectorStatus       string                `gorm:"column:vector_status;type:varchar(32);not null;default:'pending';comment:向量状态 pending/indexed/failed" json:"vectorStatus"`
	VectorID           string                `gorm:"column:vector_id;type:varchar(255);default:'';comment:Milvus 中的向量ID" json:"vectorID"`
	Metadata           datatypes.JSONMap     `gorm:"column:metadata;default:'{}';comment:扩展元数据" json:"metadata"`
	IndexingError      string                `gorm:"column:indexing_error;type:text;default:'';comment:索引错误信息" json:"indexingError"`
	LastIndexedAt      int64                 `gorm:"column:last_indexed_at;type:bigint;default:0;comment:最后索引时间" json:"lastIndexedAt"`

	// 关联集合
	Collection *KnowledgeCollection `gorm:"foreignKey:CollectionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"collection"`
}

func (*KnowledgeDocument) TableName() string {
	return "knowledge_document"
}
