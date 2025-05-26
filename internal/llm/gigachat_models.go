package llm

type gigaChatAccessToken struct {
	Token     string `json:"access_token"`
	ExpiresAt int64  `json:"expires_at"`
}

type gigaChatClient struct {
	ApiBase     string
	Model       string
	EmbedModel  string
	accessToken gigaChatAccessToken
}

type gigaChatEmbedRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type gigaChatEmbedResponse struct {
	Data   []gigachatEmbedData `json:"data"`
	Model  string              `json:"model"`
	Object string              `json:"object"`
}

type gigachatEmbedData struct {
	Object    string    `json:"object"`
	Embedding []float32 `json:"embedding"`
}

type gigaChatRequest struct {
	ChatRequest
	Functions string `json:"functions,omitempty"`
}

type GigaChatResponse struct {
	Choices []Choice   `json:"choices"`
	Created int64      `json:"created"`
	Model   string     `json:"model"`
	Usage   UsageStats `json:"usage"`
	Object  string     `json:"object"`
}

type Choice struct {
	Message          ResponseMessage   `json:"message"`
	Index            int               `json:"index"`
	FinishReason     string            `json:"finish_reason"` // stop, length, function_call, blacklist, error
	FunctionCall     *GigaFunctionCall `json:"function_call,omitempty"`
	FunctionsStateID string            `json:"functions_state_id,omitempty"` // UUIDv4
}

type ResponseMessage struct {
	Role    string `json:"role"`              // assistant, function_in_progress
	Content string `json:"content"`           // Текст или статус выполнения
	Created *int64 `json:"created,omitempty"` // Только для function_in_progress
	Name    string `json:"name,omitempty"`    // Название функции
}

type GigaFunctionCall struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

type UsageStats struct {
	PromptTokens          int `json:"prompt_tokens"`
	CompletionTokens      int `json:"completion_tokens"`
	PrecachedPromptTokens int `json:"precached_prompt_tokens"`
	TotalTokens           int `json:"total_tokens"`
}
