package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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

func NewLlamaClient() *LlamaClient {
	apiBase := os.Getenv("ASAI_LLM_API_BASE")
	model := os.Getenv("ASAI_LLM_MODEL")

	if apiBase == "" {
		apiBase = "http://localhost:11434"
	}
	if model == "" {
		model = "llama3.2:1b"
	}

	return &LlamaClient{
		ApiBase: strings.TrimRight(apiBase, "/"),
		Model:   model,
	}
}

func (c *LlamaClient) Generate(prompt string) (string, error) {
	url := fmt.Sprintf("%s/api/generate", c.ApiBase)

	reqBody := llamaRequest{
		Model:  c.Model,
		Prompt: prompt,
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
