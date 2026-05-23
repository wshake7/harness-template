package agent

import (
	"admin/internal/ai/agent/knowledge_pipeline"
	"admin/internal/ai/loader"
	"admin/internal/config"
	"admin/internal/domains"
	"admin/internal/services/milvusc"
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/compose"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

func TestKnowledgeRun(t *testing.T) {
	ctx := context.Background()
	conf := config.Init("../../../etc/config.local.yaml")
	cli, err := milvusc.New(ctx, conf.Milvus)
	if err != nil {
		t.Fatalf("init milvus client failed: %v", err)
	}

	err = cli.DropCollection(ctx, milvusclient.NewDropCollectionOption(domains.MilvusCollectionName))
	if err != nil {
		fmt.Printf("[warn] drop collection failed (may not exist): %v\n", err)
	} else {
		fmt.Printf("[info] dropped collection: %s\n", domains.MilvusCollectionName)
	}

	r, err := knowledge_pipeline.BuildKnowledgeIndexing(ctx, conf)
	if err != nil {
		t.Fatalf("build knowledge indexing failed: %v", err)
	}

	err = filepath.WalkDir(domains.MilvusFileDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walk dir failed: %w", err)
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".md") {
			fmt.Printf("[skip] not a markdown file: %s\n", path)
			return nil
		}

		fmt.Printf("[start] indexing file: %s\n", path)
		// 删除biz数据metadata中_source一样的数据
		ldr, err := loader.New(ctx)
		if err != nil {
			return err
		}
		docs, err := ldr.Load(ctx, document.Source{URI: path})
		if err != nil {
			return err
		}
		expr := fmt.Sprintf(`metadata["_source"] == "%s"`, docs[0].MetaData["_source"])
		queryResult, err := cli.Query(ctx, milvusclient.NewQueryOption(domains.MilvusCollectionName).
			WithFilter(expr).
			WithOutputFields("id"))
		if err != nil {
			return err
		} else if queryResult.Len() > 0 {
			idColumn := queryResult.GetColumn("id")
			var idsToDelete []string
			for i := 0; i < idColumn.Len(); i++ {
				id, err := idColumn.GetAsString(i)
				if err == nil {
					idsToDelete = append(idsToDelete, id)
				}
			}
			if len(idsToDelete) > 0 {
				_, err = cli.Delete(ctx, milvusclient.NewDeleteOption(domains.MilvusCollectionName).
					WithStringIDs("id", idsToDelete))
				if err != nil {
					fmt.Printf("[warn] delete existing data failed: %v\n", err)
				} else {
					fmt.Printf("[info] deleted %d existing records with _source: %s\n", len(idsToDelete), docs[0].MetaData["_source"])
				}
			}
		}
		// 重新构建
		ids, err := r.Invoke(ctx, document.Source{URI: path}, compose.WithCallbacks())
		if err != nil {
			return fmt.Errorf("invoke index graph failed: %w", err)
		}
		fmt.Printf("[done] indexing file: %s, len of parts: %d，%s\n", path, len(ids), ids)
		return nil
	})
	if err != nil {
		t.Fatalf("walk dir failed: %v", err)
	}
}
