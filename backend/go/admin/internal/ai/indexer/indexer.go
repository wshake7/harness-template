package indexer

import (
	"admin/internal/ai/embedder"
	"admin/internal/config"
	"admin/internal/services/milvusc"
	"context"
	"fmt"
	"strconv"
	"strings"

	milvus2 "github.com/cloudwego/eino-ext/components/indexer/milvus2"
	einoindexer "github.com/cloudwego/eino/components/indexer"
	"github.com/milvus-io/milvus/client/v2/entity"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

const defaultCollectionName = "biz"

func New(ctx context.Context, conf *config.Config) (einoindexer.Indexer, error) {
	if !conf.Milvus.Enabled {
		return nil, fmt.Errorf("milvus is disabled")
	}

	cli, err := getMilvusClient(ctx, conf.Milvus)
	if err != nil {
		return nil, err
	}
	eb, err := embedder.New(ctx, conf.AI.Embedding)
	if err != nil {
		return nil, err
	}

	knowledge := conf.AI.Knowledge
	collection := knowledge.Collection
	if collection == "" {
		collection = defaultCollectionName
	}

	return milvus2.NewIndexer(ctx, &milvus2.IndexerConfig{
		Client:      cli,
		Collection:  collection,
		Description: knowledge.CollectionDesc,
		Vector: &milvus2.VectorConfig{
			// Use the configured embedding dimension directly so Milvus schema stays aligned
			// with the active embedding model.
			Dimension:    int64(conf.AI.Embedding.Dimensions),
			MetricType:   metricType(conf.AI.Knowledge.IndexMetricType),
			IndexBuilder: indexBuilder(conf.AI.Knowledge.IndexType),
			VectorField:  "vector",
		},
		FieldParams: fieldParams(conf.AI.Knowledge),
		Embedding:   eb,
	})
}

func getMilvusClient(ctx context.Context, conf config.MilvusConfig) (*milvusclient.Client, error) {
	if milvusc.Client != nil {
		return milvusc.Client, nil
	}
	return milvusc.New(ctx, conf)
}

func fieldParams(conf config.AIKnowledgeConfig) map[string]map[string]string {
	return map[string]map[string]string{
		"id": {
			entity.TypeParamMaxLength: strconv.FormatUint(uint64(conf.IDMaxLength), 10),
		},
		"content": {
			entity.TypeParamMaxLength: strconv.FormatUint(uint64(conf.ContentMaxLength), 10),
		},
	}
}

func indexBuilder(indexType string) milvus2.IndexBuilder {
	switch strings.ToLower(indexType) {
	case "hnsw":
		return milvus2.NewHNSWIndexBuilder()
	case "ivf_flat":
		return milvus2.NewIVFFlatIndexBuilder()
	default:
		return milvus2.NewAutoIndexBuilder()
	}
}

func metricType(indexMetricType string) milvus2.MetricType {
	switch strings.ToUpper(indexMetricType) {
	case "IP":
		return milvus2.IP
	case "L2":
		return milvus2.L2
	default:
		return milvus2.COSINE
	}
}
