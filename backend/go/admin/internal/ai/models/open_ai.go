package models

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
)

// OpenAIChatModelConfig is the model-layer config for an OpenAI-compatible chat model.
// It is intentionally decoupled from the application config package.
type OpenAIChatModelConfig struct {
	Model   string
	APIKey  string
	BaseURL string
}

// NewChatModel creates an OpenAI-compatible chat model.
func NewChatModel(ctx context.Context, cfg OpenAIChatModelConfig) (cm model.ToolCallingChatModel, err error) {
	if cfg.Model == "" {
		return nil, fmt.Errorf("chat model name is empty")
	}
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("API key is empty for model %q, please set it in config", cfg.Model)
	}
	return openai.NewChatModel(ctx, &openai.ChatModelConfig{
		Model:   cfg.Model,
		APIKey:  cfg.APIKey,
		BaseURL: cfg.BaseURL,
	})
}
