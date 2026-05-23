package chat_pipeline

import (
	"admin/internal/ai/models"
	"admin/internal/config"
	"context"

	"github.com/cloudwego/eino/components/model"
)

func newChatModel(ctx context.Context, cfg config.ChatModelConfig) (cm model.ToolCallingChatModel, err error) {
	return models.NewChatModel(ctx, models.OpenAIChatModelConfig{
		Model:   cfg.Model,
		APIKey:  cfg.APIKey,
		BaseURL: cfg.BaseURL,
	})
}
