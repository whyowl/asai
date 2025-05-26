package llm

import (
	"asai/internal/config"
	"asai/internal/tools"
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
	"strings"
	"time"
)

func NewGigaChatClient() *gigaChatClient {
	return &gigaChatClient{
		ApiBase:     strings.TrimRight(config.AppConfig.LLM.GigaChat.ClientUrl, "/"),
		Model:       config.AppConfig.LLM.GigaChat.Model,
		accessToken: gigaChatAccessToken{Token: "", ExpiresAt: 0},
	}
}

func (c *gigaChatClient) Generate(messages []Message, tools []tools.Tool) ([]Message, error) {

	if c.accessToken.ExpiresAt <= time.Now().Unix() {
		err := c.accessRequest(config.AppConfig.LLM.GigaChat.Secret)
		if err != nil {
			return []Message{}, fmt.Errorf("error request access: %s", err)
		}
	}

	url := fmt.Sprintf("%s/api/v1/chat/completions", c.ApiBase)
	reqBody := gigaChatRequest{
		ChatRequest: ChatRequest{
			Model:    c.Model,
			Messages: messages,
		},
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return []Message{}, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return []Message{}, err
	}
	req.Header.Set("Authorization", "Bearer "+c.accessToken.Token)
	req.Header.Set("Content-Type", "application/json")

	client := newHTTPClientWithCert(config.AppConfig.LLM.GigaChat.Certificate)

	resp, err := client.Do(req)
	if err != nil {
		return []Message{}, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return []Message{}, err
	}
	var respObj GigaChatResponse
	err = json.Unmarshal(respBytes, &respObj)
	if err != nil {
		return []Message{}, fmt.Errorf("error decoding response: %v, body: %s", err, string(respBytes))
	}
	if len(respObj.Choices) == 0 {
		return []Message{}, fmt.Errorf("error empty response: %v, body: %s", err, string(respBytes))
	}
	return []Message{{Role: respObj.Choices[0].Message.Role, Content: respObj.Choices[0].Message.Content}}, nil
}

func (c *gigaChatClient) Embed(input string) ([]float32, error) {

	if c.accessToken.ExpiresAt <= time.Now().Unix() {
		err := c.accessRequest(config.AppConfig.LLM.GigaChat.Secret)
		if err != nil {
			return []float32{}, fmt.Errorf("error request access: %s", err)
		}
	}

	url := fmt.Sprintf("%s/api/v1/embeddings", c.ApiBase)
	reqBody := gigaChatEmbedRequest{
		Model: c.EmbedModel,
		Input: input,
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return []float32{}, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return []float32{}, err
	}
	req.Header.Set("Authorization", "Bearer "+c.accessToken.Token)
	req.Header.Set("Content-Type", "application/json")

	client := newHTTPClientWithCert(config.AppConfig.LLM.GigaChat.Certificate)

	resp, err := client.Do(req)
	if err != nil {
		return []float32{}, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return []float32{}, err
	}
	var respObj gigaChatEmbedResponse
	err = json.Unmarshal(respBytes, &respObj)
	if err != nil {
		return []float32{}, fmt.Errorf("error decoding response: %v, body: %s", err, string(respBytes))
	}
	if len(respObj.Data) == 0 {
		return []float32{}, fmt.Errorf("error empty response: %v, body: %s", err, string(respBytes))
	}
	return respObj.Data[0].Embedding, nil
}

func (c *gigaChatClient) accessRequest(clientSecret string) error {

	rqUID := uuid.New().String()
	form := url.Values{}
	form.Set("scope", config.AppConfig.LLM.GigaChat.Scope)

	req, err := http.NewRequest("POST", config.AppConfig.LLM.GigaChat.TokenUrl, strings.NewReader(form.Encode()))
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
