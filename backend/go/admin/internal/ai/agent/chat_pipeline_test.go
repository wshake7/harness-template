package agent

import (
	"admin/internal/ai/agent/chat_pipeline"
	"admin/internal/ai/mem"
	"admin/internal/config"
	"admin/internal/services/milvusc"
	"context"
	"fmt"
	"testing"

	"github.com/cloudwego/eino/schema"
)

func TestChatRun(t *testing.T) {
	ctx := context.Background()
	id := "111"
	userMessage := &chat_pipeline.UserMessage{
		ID:    id,
		Query: "你好",
	}
	conf := config.Init("../../../etc/config.local.yaml")
	_, err := milvusc.New(ctx, conf.Milvus)
	modelConf, err := config.ResolveModel(conf.AI.Models, "deepseek-v4-flash")
	if err != nil {
		panic(err)
	}
	runner, err := chat_pipeline.BuildChatAgent(ctx, conf.AI.Embedding, modelConf, conf.AI.McpURL, conf.AI.Knowledge.IndexMetricType)
	if err != nil {
		panic(err)
	}
	// 第一次对话
	out, err := runner.Invoke(ctx, userMessage)
	if err != nil {
		panic(err)
	}
	answer := out.Content
	fmt.Println("Q: 你好")
	fmt.Println("A:", answer)
	memory := mem.NewInMemMemory(mem.InMemConfig{MaxWindowSize: 3000})
	session := memory.GetSession(id)
	err = session.SetMessages(ctx, schema.UserMessage("你好"))
	if err != nil {
		panic(err)
	}
	err = session.SetMessages(ctx, schema.UserMessage(out.Content))
	if err != nil {
		panic(err)
	}
	// 第二次对话
	messages, err := session.GetMessages(ctx)
	if err != nil {
		panic(err)
	}
	userMessage = &chat_pipeline.UserMessage{
		ID:      id,
		Query:   "现在是几点",
		History: messages,
	}
	out, err = runner.Invoke(ctx, userMessage)
	if err != nil {
		panic(err)
	}
	answer = out.Content
	fmt.Println("----------------")
	fmt.Println("Q: 现在是几点")
	fmt.Println("A:", answer)
}
