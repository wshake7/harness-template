package retriever

import (
	"admin/internal/ai/embedder"
	"admin/internal/config"
	"admin/internal/domains"
	"admin/internal/services/milvusc"
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/eino-ext/components/retriever/milvus2"
	"github.com/cloudwego/eino-ext/components/retriever/milvus2/search_mode"
)

func NewMilvusRetriever(ctx context.Context, embConf config.AIEmbeddingConfig, metricType string) (rtr *milvus2.Retriever, err error) {
	cli := milvusc.Client
	if cli == nil {
		return nil, fmt.Errorf("milvus client not initialized")
	}
	eb, err := embedder.New(ctx, embConf)
	if err != nil {
		return nil, err
	}
	r, err := milvus2.NewRetriever(ctx, &milvus2.RetrieverConfig{
		Client:      cli,
		Collection:  domains.MilvusCollectionName,
		VectorField: "vector",
		OutputFields: []string{
			"id",
			"content",
			"metadata",
		},
		TopK:       1,
		Embedding:  eb,
		SearchMode: search_mode.NewApproximate(metricTypeFromString(metricType)),
	})
	if err != nil {
		return nil, fmt.Errorf("create milvus retriever: %w", err)
	}
	return r, nil
}

func metricTypeFromString(s string) milvus2.MetricType {
	switch strings.ToUpper(s) {
	case "IP":
		return milvus2.IP
	case "L2":
		return milvus2.L2
	default:
		return milvus2.COSINE
	}
}
