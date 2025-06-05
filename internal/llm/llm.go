package llm

import (
	"asai/internal/shared"
	"asai/internal/tools"
	"context"
)

var Providers = map[string]LLM{}

type ChatRequest struct {
	Model    string           `json:"model"`
	Messages []shared.Message `json:"messages"`
}

type LLM interface {
	Generate(ctx context.Context, prompt []shared.Message, tools []tools.Function, userId int64) ([]shared.Message, error)
	Embed(ctx context.Context, input string) ([]float32, error)
}
