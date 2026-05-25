package knowledge_pipeline

import (
	"admin/internal/config"
	"context"

	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/compose"
)

type BuildOptions struct {
	SplitterConfig  *SplitterOptions
	Metadata        map[string]any
	IndexerOverride *config.Config
}

func BuildKnowledgeIndexing(ctx context.Context, conf *config.Config) (r compose.Runnable[document.Source, []string], err error) {
	return BuildKnowledgeIndexingWithOptions(ctx, conf, nil)
}

func BuildKnowledgeIndexingWithOptions(ctx context.Context, conf *config.Config, opts *BuildOptions) (r compose.Runnable[document.Source, []string], err error) {
	const (
		FileLoader       = "FileLoader"
		MarkdownSplitter = "MarkdownSplitter"
		MetadataEnricher = "MetadataEnricher"
		MilvusIndexer    = "MilvusIndexer"
	)
	g := compose.NewGraph[document.Source, []string]()
	fileLoaderKeyOfLoader, err := newLoader(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddLoaderNode(FileLoader, fileLoaderKeyOfLoader)
	markdownSplitterKeyOfDocumentTransformer, err := newDocumentTransformer(ctx, splitterOptions(opts))
	if err != nil {
		return nil, err
	}
	_ = g.AddDocumentTransformerNode(MarkdownSplitter, markdownSplitterKeyOfDocumentTransformer)
	if len(metadataOptions(opts)) > 0 {
		enricher, metaErr := newMetadataTransformer(metadataOptions(opts))
		if metaErr != nil {
			return nil, metaErr
		}
		_ = g.AddDocumentTransformerNode(MetadataEnricher, enricher)
	}
	indexerKeyOfIndexer, err := newIndexer(ctx, confOverride(conf, opts))
	if err != nil {
		return nil, err
	}
	_ = g.AddIndexerNode(MilvusIndexer, indexerKeyOfIndexer)
	_ = g.AddEdge(compose.START, FileLoader)
	_ = g.AddEdge(MilvusIndexer, compose.END)
	_ = g.AddEdge(FileLoader, MarkdownSplitter)
	if len(metadataOptions(opts)) > 0 {
		_ = g.AddEdge(MarkdownSplitter, MetadataEnricher)
		_ = g.AddEdge(MetadataEnricher, MilvusIndexer)
	} else {
		_ = g.AddEdge(MarkdownSplitter, MilvusIndexer)
	}
	r, err = g.Compile(ctx, compose.WithGraphName("KnowledgeIndexing"), compose.WithNodeTriggerMode(compose.AnyPredecessor))
	if err != nil {
		return nil, err
	}
	return r, err
}

func splitterOptions(opts *BuildOptions) *SplitterOptions {
	if opts == nil {
		return nil
	}
	return opts.SplitterConfig
}

func metadataOptions(opts *BuildOptions) map[string]any {
	if opts == nil {
		return nil
	}
	return opts.Metadata
}

func confOverride(conf *config.Config, opts *BuildOptions) *config.Config {
	if opts == nil || opts.IndexerOverride == nil {
		return conf
	}
	return opts.IndexerOverride
}
