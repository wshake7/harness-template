package embedder

import (
	"admin/internal/config"
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/embedding/dashscope"
	"github.com/cloudwego/eino/components/embedding"
)

func New(ctx context.Context, conf config.AIEmbeddingConfig) (embedding.Embedder, error) {
	if conf.Provider != "" && conf.Provider != "dashscope" {
		return nil, fmt.Errorf("unsupported embedding provider: %s", conf.Provider)
	}
	if conf.APIKey == "" {
		return nil, fmt.Errorf("ai embedding api key is empty")
	}
	if conf.Model == "" {
		return nil, fmt.Errorf("ai embedding model is empty")
	}

	cfg := &dashscope.EmbeddingConfig{
		APIKey: conf.APIKey,
		Model:  conf.Model,
	}
	if conf.Dimensions > 0 {
		// The config keeps dimensions as an unsigned count, while the SDK expects *int.
		dim := int(conf.Dimensions)
		cfg.Dimensions = &dim
	}

	eb, err := dashscope.NewEmbedder(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("create dashscope embedder: %w", err)
	}
	return eb, nil
}
