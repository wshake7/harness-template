package agent

import (
	"admin/internal/ai/models"
	"admin/internal/ai/tools"
	"admin/internal/config"
	"context"
	"fmt"
	"testing"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func TestToolCall(t *testing.T) {
	ctx := context.Background()
	conf := config.Init("../../../etc/config.local.yaml")

	if len(conf.AI.Models) == 0 {
		t.Fatal("no AI models configured")
	}

	chatModel, err := models.NewChatModel(ctx, models.OpenAIChatModelConfig{
		Model:   conf.AI.Models[0].Model,
		APIKey:  conf.AI.Models[0].APIKey,
		BaseURL: conf.AI.Models[0].BaseURL,
	})
	if err != nil {
		t.Fatalf("create chat model: %v", err)
	}

	// 获取工具信息, 用于绑定到 ChatModel
	mcpTool, _ := tools.GetLogMcpTool(conf.AI.McpURL)
	toolList := make([]tool.BaseTool, 0, len(mcpTool)+1)
	toolList = append(toolList, mcpTool...)
	toolList = append(toolList, tools.NewGetCurrentTimeTool())

	toolInfos := make([]*schema.ToolInfo, 0, len(toolList))
	for _, bt := range toolList {
		info, err := bt.Info(ctx)
		if err != nil {
			t.Fatalf("get tool info: %v", err)
		}
		toolInfos = append(toolInfos, info)
	}

	// 将 tools 绑定到 ChatModel
	chatModel, err = chatModel.WithTools(toolInfos)
	if err != nil {
		t.Fatalf("bind tools: %v", err)
	}

	// 创建一个完整的处理链
	chain := compose.NewChain[[]*schema.Message, *schema.Message]()
	chain.AppendChatModel(chatModel, compose.WithNodeName("chat_model"))

	// 编译并运行 chain
	agent, err := chain.Compile(ctx)
	if err != nil {
		t.Fatalf("compile chain: %v", err)
	}

	// 运行示例
	resp, err := agent.Invoke(ctx, []*schema.Message{
		{
			Role:    schema.User,
			Content: "告诉我你有哪些工具可以使用",
		},
	})
	if err != nil {
		t.Fatalf("invoke agent: %v", err)
	}
	fmt.Println(resp.Content)
}
