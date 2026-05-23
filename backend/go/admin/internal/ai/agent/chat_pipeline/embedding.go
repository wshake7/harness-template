package chat_pipeline

import (
	"admin/internal/ai/embedder"
	"admin/internal/config"
	"context"

	"github.com/cloudwego/eino/components/embedding"
)

func newEmbedding(ctx context.Context, conf config.AIEmbeddingConfig) (eb embedding.Embedder, err error) {
	return embedder.New(ctx, conf)
}
