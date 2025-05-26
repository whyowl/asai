package llm

import (
	"asai/internal/config"
	"asai/internal/tools"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

func (c *OllamaClient) Generate(messages []Message, functionsTools []tools.Tool) ([]Message, error) {
	url := fmt.Sprintf("%s/api/chat", c.ApiBase)
	reqBody := ollamaRequest{
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

	var respObj ollamaResponse
	err = json.Unmarshal(respBytes, &respObj)
	if err != nil {
		return []Message{}, fmt.Errorf("error decoding response: %v, body: %s", err, string(respBytes))
	}

	var functionsResponse string
	for _, f := range respObj.Message.ToolCalls {
		functionResponse, err := tools.CallFunctionsByModel(f.Function.Name, f.Function.Arguments)
		if err != nil {
			return []Message{}, err //кривой код
		}
		functionsResponse += "\n\n" + functionResponse
	}
	if functionsResponse != "" {
		generate, err := c.Generate(append(messages, Message{Role: "tool", Content: functionsResponse}), functionsTools)
		if err != nil {
			return []Message{}, err
		}
		return generate, nil
	}
	return []Message{Message{Content: RemoveThinkTags(respObj.Message.Content), Role: respObj.Message.Role}}, nil
}

func (c *OllamaClient) Embed(input string) ([]float32, error) {
	url := fmt.Sprintf("%s/api/embed", c.ApiBase)
	reqBody := ollamaEmbedRequest{Model: c.EmbedModel, Input: input}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return []float32{}, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return []float32{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return []float32{}, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return []float32{}, err
	}

	var respObj ollamaEmbedResponse
	err = json.Unmarshal(respBytes, &respObj)
	if err != nil {
		return []float32{}, fmt.Errorf("error decoding response: %v, body: %s", err, string(respBytes))
	}
	return respObj.Embeddings[0], nil // need check error from server
}

func RemoveThinkTags(s string) string {
	re := regexp.MustCompile(`(?s)<think>.*?</think>`)
	return re.ReplaceAllString(s, "")
}
