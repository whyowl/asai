package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type UniversalLLMClient struct {
	ApiBase string
	ApiKey  string
	Model   string
}

func NewUniversalLLMClient() *UniversalLLMClient {
	return &UniversalLLMClient{
		ApiBase: os.Getenv("ASAI_LLM_API_BASE"),
		ApiKey:  os.Getenv("ASAI_LLM_API_KEY"),
		Model:   os.Getenv("ASAI_LLM_MODEL"),
	}
}

func (c *UniversalLLMClient) Generate(prompt string) (string, error) {
	url := fmt.Sprintf("%s/chat/completions", c.ApiBase)

	body := map[string]interface{}{
		"model": c.Model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}

	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+c.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("unexpected response format")
	}

	message := choices[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
	return message, nil
}
