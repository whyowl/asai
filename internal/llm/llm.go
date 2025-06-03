package llm

import (
	"asai/internal/shared"
	"asai/internal/tools"
)

var Providers = map[string]LLM{}

type ChatRequest struct {
	Model    string           `json:"model"`
	Messages []shared.Message `json:"messages"`
}

type LLM interface {
	Generate(prompt []shared.Message, tools []tools.Function, userId int64) ([]shared.Message, error)
	Embed(input string) ([]float32, error)
}
