package appsvc

import (
	"admin/internal/config"
	"admin/internal/services/orm/models"
	"context"
)

//go:generate mockgen -source=knowledge_document_indexer.go -destination=../mock/mock_knowledge_document_indexer.go -package=mock -typed

type KnowledgeDocumentIndexer interface {
	ImportFile(ctx context.Context, input ImportKnowledgeFileInput) (*models.KnowledgeDocument, error)
	IndexDocument(ctx context.Context, documentID uint64) error
	DeleteDocumentVectors(ctx context.Context, documentID uint64) error
	EnsureCollection(ctx context.Context, collection *models.KnowledgeCollection) error
	DropCollection(ctx context.Context, collectionName string) error
}

type ImportKnowledgeFileInput struct {
	CollectionID uint64
	FileAssetID  uint64
	Title        string
	ContentType  string
	Remark       string
	OperatorID   uint64
}

type KnowledgeIndexRunner interface {
	Index(ctx context.Context, conf *config.Config, collection *models.KnowledgeCollection, sourcePath string, metadata map[string]any) ([]string, error)
}

type VectorStoreCleaner interface {
	DeleteByDocument(ctx context.Context, collectionName string, documentDBID uint64, source string) error
}
