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
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type llamaResponse struct {
	Model      string  `json:"model"`
	CreatedAt  string  `json:"created_at"` //уточнить тип данных
	Message    Message `json:"message"`
	DoneReason string  `json:"done_reason"`
	Done       bool    `json:"done"`
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

func (c *LlamaClient) Generate(messages []Message) (Message, error) {
	url := fmt.Sprintf("%s/api/chat", c.ApiBase)
	reqBody := llamaRequest{
		Model:    c.Model,
		Messages: messages,
		Stream:   false,
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return Message{}, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return Message{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Message{}, err
	}
	defer resp.Body.Close()

	if os.Getenv("DEBUG") == "true" {
		fmt.Printf("response from url: %s\n", resp)
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return Message{}, err
	}

	var respObj llamaResponse
	err = json.Unmarshal(respBytes, &respObj)
	if err != nil {
		return Message{}, fmt.Errorf("error decoding response: %v, body: %s", err, string(respBytes))
	}

	if os.Getenv("DEBUG") == "true" {
		fmt.Printf("response json: %s\n", respObj)
	}

	return respObj.Message, nil
}
