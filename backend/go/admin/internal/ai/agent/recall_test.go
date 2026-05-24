package agent

import (
	retriever2 "admin/internal/ai/retriever"
	"admin/internal/config"
	"admin/internal/services/milvusc"
	"context"
	"fmt"
	"testing"
)

func TestReCall(t *testing.T) {
	ctx := context.Background()
	conf := config.Init("../../../etc/config.local.yaml")
	_, err := milvusc.New(ctx, conf.Milvus)
	if err != nil {
		t.Fatalf("init milvus client failed: %v", err)
	}
	r, err := retriever2.NewMilvusRetriever(ctx, conf.AI.Embedding, "")
	if err != nil {
		panic(err)
	}
	query := "服务下线是什么原因"
	docs, err := r.Retrieve(ctx, query)
	if err != nil {
		panic(err)
	}
	fmt.Println("Q：", query)
	for _, doc := range docs {
		fmt.Println("A：", doc.Content)
	}
	fmt.Println("Done", len(docs))
}
