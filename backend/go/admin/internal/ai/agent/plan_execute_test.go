package agent

import (
	"admin/internal/ai/agent/plan_execute_replan"
	"admin/internal/ai/models"
	"admin/internal/ai/tools"
	"admin/internal/config"
	"admin/internal/services/milvusc"
	"context"
	"fmt"
	"testing"

	"github.com/cloudwego/eino/components/tool"
)

func TestPlanExecute(t *testing.T) {
	ctx := context.Background()
	conf := config.Init("../../../etc/config.local.yaml")

	if len(conf.AI.Models) == 0 {
		t.Fatal("no AI models configured")
	}

	_, err := milvusc.New(ctx, conf.Milvus)
	if err != nil {
		t.Fatalf("init milvus client failed: %v", err)
	}

	// Planner and Replanner need tool calling; disable thinking mode
	// so models like deepseek-reasoner / deepseek-v3.1 don't reject tool_choice.
	planModel, err := models.NewChatModel(ctx, models.OpenAIChatModelConfig{
		Model:   conf.AI.Models[0].Model,
		APIKey:  conf.AI.Models[0].APIKey,
		BaseURL: conf.AI.Models[0].BaseURL,
		ExtraFields: map[string]any{
			"thinking": map[string]string{"type": "disabled"},
		},
	})
	if err != nil {
		t.Fatalf("create plan model: %v", err)
	}
	execModel, err := models.NewChatModel(ctx, models.OpenAIChatModelConfig{
		Model:   conf.AI.Models[0].Model,
		APIKey:  conf.AI.Models[0].APIKey,
		BaseURL: conf.AI.Models[0].BaseURL,
	})
	if err != nil {
		t.Fatalf("create exec model: %v", err)
	}
	replanModel, err := models.NewChatModel(ctx, models.OpenAIChatModelConfig{
		Model:   conf.AI.Models[0].Model,
		APIKey:  conf.AI.Models[0].APIKey,
		BaseURL: conf.AI.Models[0].BaseURL,
		ExtraFields: map[string]any{
			"thinking": map[string]string{"type": "disabled"},
		},
	})
	if err != nil {
		t.Fatalf("create replan model: %v", err)
	}

	var toolList []tool.BaseTool
	mcpTool, _ := tools.GetLogMcpTool(conf.AI.McpURL)
	toolList = append(toolList, mcpTool...)
	toolList = append(toolList, tools.NewPrometheusAlertsQueryTool())
	toolList = append(toolList, tools.NewQueryInternalDocsTool(conf.AI.Embedding, conf.AI.Knowledge.IndexMetricType))
	toolList = append(toolList, tools.NewGetCurrentTimeTool())
	toolList = append(toolList, tools.NewDBCrudTool())

	answer, detail, err := plan_execute_replan.BuildPlanAgent(ctx, planModel, execModel, replanModel, toolList, "查询最近的错误日志并分析原因")
	if err != nil {
		t.Fatalf("build plan agent: %v", err)
	}
	fmt.Println("Answer:", answer)
	for i, d := range detail {
		fmt.Printf("Detail[%d]: %s\n", i, d)
	}
}
