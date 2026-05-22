package knowledge_pipeline

import (
	"admin/internal/ai/indexer"
	"admin/internal/config"
	"context"

	einoindexer "github.com/cloudwego/eino/components/indexer"
)

// newIndexer component initialization function of node 'Indexer' in graph 'KnowledgeIndexing'
func newIndexer(ctx context.Context, conf *config.Config) (idr einoindexer.Indexer, err error) {
	return indexer.New(ctx, conf)
}
