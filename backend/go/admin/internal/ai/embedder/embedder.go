package embedder

import (
	"admin/internal/config"
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino/components/embedding"
)

func New(ctx context.Context, conf config.AIEmbeddingConfig) (embedding.Embedder, error) {
	if conf.Provider != "" && conf.Provider != "ark" {
		return nil, fmt.Errorf("unsupported embedding provider: %s", conf.Provider)
	}
	if conf.APIKey == "" {
		return nil, fmt.Errorf("ai embedding api key is empty")
	}
	if conf.Model == "" {
		return nil, fmt.Errorf("ai embedding model is empty")
	}

	cfg := &ark.EmbeddingConfig{
		APIKey: conf.APIKey,
		Model:  conf.Model,
	}

	if conf.Model == "doubao-embedding-vision-251215" || conf.Model == "doubao-embedding-vision-241215" {
		apiType := ark.APITypeMultiModal
		cfg.APIType = &apiType
	}

	eb, err := ark.NewEmbedder(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("create ark embedder: %w", err)
	}
	return eb, nil
}
