package chat_pipeline

import (
	retriever2 "admin/internal/ai/retriever"
	"admin/internal/config"
	"context"

	"github.com/cloudwego/eino/components/retriever"
)

// newRetriever component initialization function of node 'MilvusRetriever' in graph 'EinoAgent'
func newRetriever(ctx context.Context, embConf config.AIEmbeddingConfig, metricType string) (rtr retriever.Retriever, err error) {
	return retriever2.NewMilvusRetriever(ctx, embConf, metricType)
}
