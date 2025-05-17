package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type LlamaClient struct {
	ApiBase string
	Model   string
}

type llamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type llamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func NewLlamaClient(uri string, model string) *LlamaClient {

	if uri == "" {
		uri = "http://localhost:11434"
	}
	if model == "" {
		model = "llama3.1:8b"
	}

	return &LlamaClient{
		ApiBase: strings.TrimRight(uri, "/"),
		Model:   model,
	}
}

func (c *LlamaClient) Generate(prompt []Message) (string, error) {
	url := fmt.Sprintf("%s/api/generate", c.ApiBase)
	fmt.Println(prompt)
	reqBody := llamaRequest{
		Model:  c.Model,
		Prompt: FormatMessagesForPrompt(prompt),
		Stream: false,
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var respObj llamaResponse
	err = json.Unmarshal(respBytes, &respObj)
	if err != nil {
		return "", fmt.Errorf("error decoding response: %v, body: %s", err, string(respBytes))
	}

	return respObj.Response, nil
}

func FormatMessagesForPrompt(messages []Message) string {
	var sb strings.Builder

	for _, msg := range messages {
		switch msg.Role {
		case "system":
			sb.WriteString("### System:\n")
		case "user":
			sb.WriteString("### User:\n")
		case "assistant":
			sb.WriteString("### Assistant:\n")
		}
		sb.WriteString(msg.Content)
		sb.WriteString("\n\n")
	}

	return sb.String()
}
