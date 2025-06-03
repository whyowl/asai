package llm

import "asai/internal/tools"

type gigachatAccessToken struct {
	Token     string `json:"access_token"`
	ExpiresAt int64  `json:"expires_at"`
}

type gigachatClient struct {
	ApiBase     string
	Model       string
	EmbedModel  string
	accessToken gigachatAccessToken
}

type gigachatEmbedRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type gigachatEmbedResponse struct {
	Data   []gigachatEmbedData `json:"data"`
	Model  string              `json:"model"`
	Object string              `json:"object"`
}

type gigachatEmbedData struct {
	Object    string    `json:"object"`
	Embedding []float32 `json:"embedding"`
}

type gigachatRequest struct {
	ChatRequest
	Functions []tools.Function `json:"functions,omitempty"`
}

type gigachatResponse struct {
	Choices []Choice   `json:"choices"`
	Created int64      `json:"created"`
	Model   string     `json:"model"`
	Usage   UsageStats `json:"usage"`
	Object  string     `json:"object"`
}

type Choice struct {
	Message          ResponseMessage `json:"message"`
	Index            int             `json:"index"`
	FinishReason     string          `json:"finish_reason"`                // stop, length, function_call, blacklist, error
	FunctionsStateID string          `json:"functions_state_id,omitempty"` // UUIDv4
}

type ResponseMessage struct {
	Role         string             `json:"role"`    // assistant, function_in_progress
	Content      string             `json:"content"` // Текст или статус выполнения
	FunctionCall tools.FunctionCall `json:"function_call,omitempty"`
	Created      *int64             `json:"created,omitempty"` // Только для function_in_progress
	Name         string             `json:"name,omitempty"`    // Название функции
}

type UsageStats struct {
	PromptTokens          int `json:"prompt_tokens"`
	CompletionTokens      int `json:"completion_tokens"`
	PrecachedPromptTokens int `json:"precached_prompt_tokens"`
	TotalTokens           int `json:"total_tokens"`
}
