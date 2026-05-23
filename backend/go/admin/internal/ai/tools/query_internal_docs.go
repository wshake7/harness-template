package tools

import (
	"admin/internal/ai/retriever"
	"admin/internal/config"
	"context"
	"encoding/json"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

type QueryInternalDocsInput struct {
	Query string `json:"query" jsonschema:"description=The query string to search in internal documentation for relevant information"`
}

func NewQueryInternalDocsTool(embConf config.AIEmbeddingConfig, metricType string) tool.InvokableTool {
	t, err := utils.InferOptionableTool(
		"query_internal_docs",
		"Search internal documentation and knowledge base for relevant information using RAG (Retrieval-Augmented Generation). Use this tool when you need to understand internal procedures, best practices, or step-by-step guides.",
		func(ctx context.Context, input *QueryInternalDocsInput, opts ...tool.Option) (output string, err error) {
			rr, err := retriever.NewMilvusRetriever(ctx, embConf, metricType)
			if err != nil {
				return "", err
			}
			resp, err := rr.Retrieve(ctx, input.Query)
			if err != nil {
				return "", err
			}
			b, _ := json.Marshal(resp)
			return string(b), nil
		})
	if err != nil {
		panic(err)
	}
	return t
}
