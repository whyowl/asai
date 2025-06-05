package llm

import (
	"asai/internal/config"
	"asai/internal/shared"
	"asai/internal/tools"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func NewOllamaClient() *OllamaClient {
	return &OllamaClient{
		ApiBase:    strings.TrimRight(config.AppConfig.LLM.Ollama.Url, "/"),
		Model:      config.AppConfig.LLM.Ollama.Model,
		EmbedModel: config.AppConfig.LLM.Ollama.EmbedModel,
	}
}

func (c *OllamaClient) Generate(ctx context.Context, messages []shared.Message, functions []tools.Function, userID int64) ([]shared.Message, error) {
	url := fmt.Sprintf("%s/api/chat", c.ApiBase)

	reqBody, err := c.buildRequestBody(messages, functions)
	if err != nil {
		return nil, fmt.Errorf("failed to build request body: %w", err)
	}

	respBytes, err := c.sendRequest(ctx, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	var respObj ollamaResponse
	if err := json.Unmarshal(respBytes, &respObj); err != nil {
		return nil, fmt.Errorf("error decoding response: %v, body: %s", err, string(respBytes))
	}

	var response []shared.Message

	functionsResponse, err := handleToolCalls(ctx, respObj.Message.ToolCalls, userID)
	if err != nil {
		log.Printf("tool call error: %v", err)
	}

	if functionsResponse != "" {
		nextMessages := append(messages, shared.Message{Role: "tool", Content: functionsResponse})
		generated, err := c.Generate(ctx, nextMessages, functions, userID)
		if err != nil {
			return nil, err
		}
		response = append(response, shared.Message{Role: "tool", Content: functionsResponse})
		response = append(response, generated...)
		return response, nil
	}

	response = append(response, shared.Message{
		Role:    respObj.Message.Role,
		Content: RemoveThinkTags(respObj.Message.Content),
	})

	return response, nil
}

func (c *OllamaClient) Embed(ctx context.Context, input string) ([]float32, error) {
	url := fmt.Sprintf("%s/api/embed", c.ApiBase)

	reqBody, err := json.Marshal(ollamaEmbedRequest{Model: c.EmbedModel, Input: input})
	if err != nil {
		return []float32{}, err
	}

	respBytes, err := c.sendRequest(ctx, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	var respObj ollamaEmbedResponse
	err = json.Unmarshal(respBytes, &respObj)
	if err != nil {
		return []float32{}, fmt.Errorf("error decoding response: %v, body: %s", err, string(respBytes))
	}

	return respObj.Embeddings[0], nil // need check error from server
}

func (c *OllamaClient) buildRequestBody(messages []shared.Message, functions []tools.Function) ([]byte, error) {
	reqBody := ollamaRequest{
		ChatRequest: ChatRequest{
			Model:    c.Model,
			Messages: messages,
		},
		Stream: false,
		Tools:  funcToTools(functions),
	}
	return json.Marshal(reqBody)
}

func (c *OllamaClient) sendRequest(ctx context.Context, url string, body []byte) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func handleToolCalls(ctx context.Context, calls []ToolCall, userID int64) (string, error) {
	var result string
	for _, f := range calls {
		res, err := tools.CallFunctionsByModel(ctx, f.Function.Name, f.Function.Arguments, userID)
		if err != nil {
			log.Printf("error with tool %s: %v", f.Function.Name, err)
			result += fmt.Sprintf("\n\n %s: error", f.Function.Name)
			continue
		}
		result += fmt.Sprintf("\n\n %s: %s", f.Function.Name, res)
	}
	return result, nil
}

func RemoveThinkTags(s string) string {
	re := regexp.MustCompile(`(?s)<think>.*?</think>`)
	return re.ReplaceAllString(s, "")
}

func funcToTools(functions []tools.Function) []ollamaTool {
	var ollamaTools []ollamaTool
	for _, f := range functions {
		ollamaTools = append(ollamaTools, ollamaTool{Type: "function", Function: f})
	}
	return ollamaTools
}
