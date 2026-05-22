package loader

import (
	"context"

	fileloader "github.com/cloudwego/eino-ext/components/document/loader/file"
	"github.com/cloudwego/eino/components/document"
)

func New(ctx context.Context) (document.Loader, error) {
	return fileloader.NewFileLoader(ctx, &fileloader.FileLoaderConfig{})
}
