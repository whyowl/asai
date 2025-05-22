package llm

import (
	"asai/internal/tools"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type LlamaClient struct {
	ApiBase string
	Model   string
}

type llamaRequest struct {
	ChatRequest
	Stream bool         `json:"stream"`
	Tools  []tools.Tool `json:"tools,omitempty"` //   data type
}

type ToolCalls struct {
	Function FunctionCall `json:"function"`
}

type FunctionCall struct {
	Name      string            `json:"name"`
	Arguments map[string]string `json:"arguments"`
}

type llamaMessageResult struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	//ToolCalls json.RawMessage `json:"tool_calls,omitempty"`
	ToolCalls []ToolCalls `json:"tool_calls,omitempty"`
}

type llamaResponse struct {
	Model      string             `json:"model"`
	CreatedAt  string             `json:"created_at"` //уточнить тип данных
	Message    llamaMessageResult `json:"message"`
	DoneReason string             `json:"done_reason"`
	Done       bool               `json:"done"`
}

func NewLlamaClient(uri string, model string) *LlamaClient {

	if uri == "" {
		uri = "http://localhost:11434"
	}
	if model == "" {
		model = "qwen3:14b"
	}

	return &LlamaClient{
		ApiBase: strings.TrimRight(uri, "/"),
		Model:   model,
	}
}

func (c *LlamaClient) Generate(messages []Message, functionsTools []tools.Tool) ([]Message, error) {
	url := fmt.Sprintf("%s/api/chat", c.ApiBase)
	reqBody := llamaRequest{
		ChatRequest: ChatRequest{
			Model:    c.Model,
			Messages: messages,
		},
		Stream: false,
		Tools:  functionsTools,
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return []Message{}, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return []Message{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return []Message{}, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return []Message{}, err
	}

	var respObj llamaResponse
	err = json.Unmarshal(respBytes, &respObj)
	if err != nil {
		return []Message{}, fmt.Errorf("error decoding response: %v, body: %s", err, string(respBytes))
	}
	if len(respObj.Message.ToolCalls) > 0 {
		var functionsResponse string
		for _, f := range respObj.Message.ToolCalls {
			functionResponse, err := tools.CallFunctionsByModel(f.Function.Name, f.Function.Arguments)
			if err != nil {
				return []Message{}, err
			}
			functionsResponse += "\n\n" + functionResponse
		}
		generate, err := c.Generate(append(messages, Message{Role: "tool", Content: functionsResponse}), functionsTools)
		if err != nil {
			return []Message{}, err
		}
		return generate, nil
	}
	return []Message{Message{Content: RemoveThinkTags(respObj.Message.Content), Role: respObj.Message.Role}}, nil
}

func RemoveThinkTags(s string) string {
	re := regexp.MustCompile(`(?s)<think>.*?</think>`)
	return re.ReplaceAllString(s, "")
}
