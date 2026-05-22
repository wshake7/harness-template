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

func New(ctx context.Context) (einoindexer.Indexer, error) {
	if !config.Conf.Milvus.Enabled {
		return nil, fmt.Errorf("milvus is disabled")
	}

	cli, err := getMilvusClient(ctx)
	if err != nil {
		return nil, err
	}
	eb, err := embedder.New(ctx)
	if err != nil {
		return nil, err
	}

	knowledge := config.Conf.AI.Knowledge
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
			Dimension:    int64(config.Conf.AI.Embedding.Dimensions),
			MetricType:   metricType(),
			IndexBuilder: indexBuilder(),
			VectorField:  "vector",
		},
		FieldParams: fieldParams(),
		Embedding:   eb,
	})
}

func getMilvusClient(ctx context.Context) (*milvusclient.Client, error) {
	if milvusc.Client != nil {
		return milvusc.Client, nil
	}
	return milvusc.New(ctx, config.Conf.Milvus)
}

func fieldParams() map[string]map[string]string {
	return map[string]map[string]string{
		"id": {
			entity.TypeParamMaxLength: strconv.FormatUint(uint64(config.Conf.AI.Knowledge.IDMaxLength), 10),
		},
		"content": {
			entity.TypeParamMaxLength: strconv.FormatUint(uint64(config.Conf.AI.Knowledge.ContentMaxLength), 10),
		},
	}
}

func indexBuilder() milvus2.IndexBuilder {
	switch strings.ToLower(config.Conf.AI.Knowledge.IndexType) {
	case "hnsw":
		return milvus2.NewHNSWIndexBuilder()
	case "ivf_flat":
		return milvus2.NewIVFFlatIndexBuilder()
	default:
		return milvus2.NewAutoIndexBuilder()
	}
}

func metricType() milvus2.MetricType {
	switch strings.ToUpper(config.Conf.AI.Knowledge.IndexMetricType) {
	case "IP":
		return milvus2.IP
	case "L2":
		return milvus2.L2
	default:
		return milvus2.COSINE
	}
}
