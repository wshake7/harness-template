package knowledge_pipeline

import (
	"context"

	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown"
	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
)

type SplitterOptions struct {
	Headers     map[string]string
	TrimHeaders bool
	IDGenerator func(ctx context.Context, originalID string, splitIndex int) string
}

// newDocumentTransformer component initialization function of node 'MarkdownSplitter' in graph 'KnowledgeIndexing'
func newDocumentTransformer(ctx context.Context, opts *SplitterOptions) (tfr document.Transformer, err error) {
	config := &markdown.HeaderConfig{
		Headers: map[string]string{
			"#": "title",
		},
		TrimHeaders: false,
		IDGenerator: func(ctx context.Context, originalID string, splitIndex int) string {
			return uuid.New().String()
		},
	}
	if opts != nil {
		if len(opts.Headers) > 0 {
			config.Headers = opts.Headers
		}
		config.TrimHeaders = opts.TrimHeaders
		if opts.IDGenerator != nil {
			config.IDGenerator = opts.IDGenerator
		}
	}
	tfr, err = markdown.NewHeaderSplitter(ctx, config)
	if err != nil {
		return nil, err
	}
	return tfr, nil
}

type metadataTransformer struct {
	metadata map[string]any
}

func newMetadataTransformer(metadata map[string]any) (document.Transformer, error) {
	return &metadataTransformer{metadata: metadata}, nil
}

func (t *metadataTransformer) Transform(ctx context.Context, docs []*schema.Document, opts ...document.TransformerOption) ([]*schema.Document, error) {
	for _, doc := range docs {
		if doc == nil {
			continue
		}
		if doc.MetaData == nil {
			doc.MetaData = map[string]any{}
		}
		for key, value := range t.metadata {
			doc.MetaData[key] = value
		}
	}
	return docs, nil
}
