package llm

import "asai/internal/tools"

var Providers = map[string]LLM{}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type LLM interface {
	Generate(prompt []Message, tools []tools.Function, userId int64) ([]Message, error)
	Embed(input string) ([]float32, error)
}
