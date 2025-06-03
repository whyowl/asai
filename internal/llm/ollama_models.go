package llm

import "asai/internal/tools"

type OllamaClient struct {
	ApiBase    string
	Model      string
	EmbedModel string
}

type ollamaEmbedRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type ollamaEmbedResponse struct {
	Model      string      `json:"model"`
	Embeddings [][]float32 `json:"embeddings"`
}

type ollamaRequest struct {
	ChatRequest
	Stream bool         `json:"stream"`
	Tools  []ollamaTool `json:"tools,omitempty"`
}

type ollamaTool struct {
	Type     string         `json:"type"` // "function"
	Function tools.Function `json:"function"`
}

type ToolCalls struct {
	Function tools.FunctionCall `json:"function"`
}

type ollamaMessageResult struct {
	Role      string      `json:"role"`
	Content   string      `json:"content"`
	ToolCalls []ToolCalls `json:"tool_calls,omitempty"`
}

type ollamaResponse struct {
	Model      string              `json:"model"`
	CreatedAt  string              `json:"created_at"` //уточнить тип данных
	Message    ollamaMessageResult `json:"message"`
	DoneReason string              `json:"done_reason"`
	Done       bool                `json:"done"`
}
