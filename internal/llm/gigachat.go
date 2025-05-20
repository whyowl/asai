package llm

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	tokenURL = "https://ngw.devices.sberbank.ru:9443/api/v2/oauth" // Заменить на фактический URL API
	scope    = "GIGACHAT_API_PERS"                                 // Или GIGACHAT_API_PERS, GIGACHAT_API_CORP
)

type gigaChatAccessToken struct {
	Token     string `json:"access_token"`
	ExpiresAt int64  `json:"expires_at"`
}

type gigaChatClient struct {
	ApiBase     string
	Model       string
	accessToken gigaChatAccessToken
}

type gigaChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type GigaChatResponse struct {
	Choices []Choice   `json:"choices"`
	Created int64      `json:"created"`
	Model   string     `json:"model"`
	Usage   UsageStats `json:"usage"`
	Object  string     `json:"object"`
}

type Choice struct {
	Message          ResponseMessage `json:"message"`
	Index            int             `json:"index"`
	FinishReason     string          `json:"finish_reason"` // stop, length, function_call, blacklist, error
	FunctionCall     *FunctionCall   `json:"function_call,omitempty"`
	FunctionsStateID string          `json:"functions_state_id,omitempty"` // UUIDv4
}

type ResponseMessage struct {
	Role    string `json:"role"`              // assistant, function_in_progress
	Content string `json:"content"`           // Текст или статус выполнения
	Created *int64 `json:"created,omitempty"` // Только для function_in_progress
	Name    string `json:"name,omitempty"`    // Название функции
}

type FunctionCall struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"` // универсальный формат
}

type UsageStats struct {
	PromptTokens          int `json:"prompt_tokens"`
	CompletionTokens      int `json:"completion_tokens"`
	PrecachedPromptTokens int `json:"precached_prompt_tokens"`
	TotalTokens           int `json:"total_tokens"`
}

func NewGigaChatClient(uri string, model string) *gigaChatClient {

	if uri == "" {
		uri = "https://gigachat.devices.sberbank.ru/"
	}
	if model == "" {
		model = "GigaChat"
	}

	return &gigaChatClient{
		ApiBase:     strings.TrimRight(uri, "/"),
		Model:       model,
		accessToken: gigaChatAccessToken{Token: "", ExpiresAt: 0},
	}
}

func (c *gigaChatClient) Generate(prompt []Message) (string, error) {

	if c.accessToken.ExpiresAt <= time.Now().Unix() {
		err := c.accessRequest(os.Getenv("GIGACHAT_CLIENT_SECRET"))
		if err != nil {
			return "", fmt.Errorf("error request access: %s", err)
		}
	}

	url := fmt.Sprintf("%s/api/v1/chat/completions", c.ApiBase)
	reqBody := gigaChatRequest{
		Model:    c.Model,
		Messages: prompt,
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.accessToken.Token)
	req.Header.Set("Content-Type", "application/json")

	client := newHTTPClientWithCert(os.Getenv("GIGACHAT_CLIENT_CERT"))

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var respObj GigaChatResponse
	err = json.Unmarshal(respBytes, &respObj)
	if err != nil {
		return "", fmt.Errorf("error decoding response: %v, body: %s", err, string(respBytes))
	}
	if len(respObj.Choices) == 0 {
		return "", fmt.Errorf("error empty response: %v, body: %s", err, string(respBytes))
	}
	return respObj.Choices[0].Message.Content, nil
}

func (c *gigaChatClient) accessRequest(clientSecret string) error {

	rqUID := uuid.New().String()
	form := url.Values{}
	form.Set("scope", scope)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Basic "+clientSecret)
	req.Header.Set("RqUID", rqUID)

	client := newHTTPClientWithCert(os.Getenv("GIGACHAT_CLIENT_CERT"))
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(respBody))
	}

	if err := json.Unmarshal(respBody, &c.accessToken); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if c.accessToken.Token == "" {
		return errors.New("access_token not found in response")
	}

	return nil
}

func newHTTPClientWithCert(certPath string) *http.Client {
	caCert, err := ioutil.ReadFile(certPath)
	if err != nil {
		log.Fatalf("failed to read cert file: %v", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		RootCAs: caCertPool,
	}

	transport := &http.Transport{TLSClientConfig: tlsConfig}
	return &http.Client{Transport: transport}
}
