package llm

import (
	"asai/internal/config"
	"asai/internal/shared"
	"asai/internal/tools"
	"bytes"
	"context"
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
	"strings"
	"time"
)

func NewGigaChatClient() *gigachatClient {
	return &gigachatClient{
		ApiBase:     strings.TrimRight(config.AppConfig.LLM.GigaChat.ClientUrl, "/"),
		Model:       config.AppConfig.LLM.GigaChat.Model,
		EmbedModel:  config.AppConfig.LLM.GigaChat.EmbedModel,
		accessToken: gigachatAccessToken{Token: "", ExpiresAt: 0},
	}
}

func (c *gigachatClient) Generate(ctx context.Context, messages []shared.Message, functions []tools.Function, userID int64) ([]shared.Message, error) {

	if err := c.ensureAccessToken(ctx); err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/chat/completions", c.ApiBase)

	reqBody, err := c.buildRequestBody(messages, functions)
	if err != nil {
		return nil, fmt.Errorf("failed to build request body: %w", err)
	}

	respBytes, err := c.sendRequest(ctx, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	var respObj gigachatResponse

	if err = json.Unmarshal(respBytes, &respObj); err != nil {
		return nil, fmt.Errorf("error decoding response: %v, body: %s", err, string(respBytes))
	}

	if len(respObj.Choices) == 0 {
		return nil, fmt.Errorf("error empty response: %v, body: %s", err, string(respBytes))
	}

	response := shared.Message{
		Content: respObj.Choices[0].Message.Content,
		Role:    respObj.Choices[0].Message.Role,
	}

	return []shared.Message{response}, nil
}

func (c *gigachatClient) buildRequestBody(messages []shared.Message, functions []tools.Function) ([]byte, error) {
	reqBody := gigachatRequest{
		ChatRequest: ChatRequest{
			Model:    c.Model,
			Messages: messages,
		},
		Functions: functions,
	}
	return json.Marshal(reqBody)
}

func (c *gigachatClient) sendRequest(ctx context.Context, url string, body []byte) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.accessToken.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := newHTTPClientWithCert(config.AppConfig.LLM.GigaChat.Certificate).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (c *gigachatClient) Embed(ctx context.Context, input string) ([]float32, error) {

	if err := c.ensureAccessToken(ctx); err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/embeddings", c.ApiBase)
	reqBody := gigachatEmbedRequest{
		Model: c.EmbedModel,
		Input: input,
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	respBytes, err := c.sendRequest(ctx, url, bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	var respObj gigachatEmbedResponse
	err = json.Unmarshal(respBytes, &respObj)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %v, body: %s", err, string(respBytes))
	}
	if len(respObj.Data) == 0 {
		return nil, fmt.Errorf("error empty response: %v, body: %s", err, string(respBytes))
	}
	return respObj.Data[0].Embedding, nil
}

func (c *gigachatClient) ensureAccessToken(ctx context.Context) error {
	if c.accessToken.ExpiresAt > time.Now().Unix() {
		return nil
	}
	return c.accessRequest(ctx, config.AppConfig.LLM.GigaChat.Secret)
}

func (c *gigachatClient) accessRequest(ctx context.Context, clientSecret string) error {

	rqUID := uuid.New().String()
	form := url.Values{}
	form.Set("scope", config.AppConfig.LLM.GigaChat.Scope)

	req, err := http.NewRequestWithContext(ctx, "POST", config.AppConfig.LLM.GigaChat.TokenUrl, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Basic "+clientSecret)
	req.Header.Set("RqUID", rqUID)

	client := newHTTPClientWithCert(config.AppConfig.LLM.GigaChat.Certificate)
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
