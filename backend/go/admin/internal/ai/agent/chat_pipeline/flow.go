package chat_pipeline

import (
	"admin/internal/ai/tools"
	"admin/internal/config"
	"context"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"go.uber.org/zap"
)

func newReactAgentLambda(ctx context.Context, embConf config.AIEmbeddingConfig, modelConf config.ChatModelConfig, mcpURL string, metricType string) (lba *compose.Lambda, err error) {
	cfg := &react.AgentConfig{
		MaxStep:            25,
		ToolReturnDirectly: map[string]struct{}{}}
	chatModelIns11, err := newChatModel(ctx, modelConf)
	if err != nil {
		return nil, err
	}
	cfg.ToolCallingModel = chatModelIns11

	mcpTool, err := tools.GetLogMcpTool(mcpURL)
	if err != nil {
		zap.L().Warn("MCP log tool not available", zap.Error(err))
	} else {
		cfg.ToolsConfig.Tools = append(cfg.ToolsConfig.Tools, mcpTool...)
	}
	cfg.ToolsConfig.Tools = append(cfg.ToolsConfig.Tools, tools.NewPrometheusAlertsQueryTool())
	cfg.ToolsConfig.Tools = append(cfg.ToolsConfig.Tools, tools.NewDBCrudTool())
	cfg.ToolsConfig.Tools = append(cfg.ToolsConfig.Tools, tools.NewGetCurrentTimeTool())
	cfg.ToolsConfig.Tools = append(cfg.ToolsConfig.Tools, tools.NewQueryInternalDocsTool(embConf, metricType))

	ins, err := react.NewAgent(ctx, cfg)
	if err != nil {
		return nil, err
	}
	lba, err = compose.AnyLambda(ins.Generate, ins.Stream, nil, nil)
	if err != nil {
		return nil, err
	}
	return lba, nil
}
