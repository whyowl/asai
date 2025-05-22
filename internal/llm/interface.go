package llm

import "asai/internal/tools"

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

//type MessageResult struct {
//	Role    string `json:"role"`
//	Content string `json:"content"`
//ToolCalls []ToolCall `json:"tool_calls,omitempty"`
//}

//type ToolCall struct {
//	Function ToolFunctionCall `json:"function"`
//}

//type ToolFunctionCall struct {
//	Name      string                 `json:"name"`
//	Arguments map[string]interface{} `json:"arguments"`
//}

type LLM interface {
	Generate(prompt []Message, tools []tools.Tool) ([]Message, error)
}
