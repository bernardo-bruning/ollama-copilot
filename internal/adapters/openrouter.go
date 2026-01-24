package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bernardo-bruning/ollama-copilot/internal/ports"
)

type OpenRouter struct {
	client  *http.Client
	token   string
	model   string
	system  string
	BaseURL string
}

func NewOpenRouter(token string, model string, system string) *OpenRouter {
	return NewOpenRouterWithBaseURL(token, model, system, "https://openrouter.ai/api/v1")
}

func NewOpenRouterWithBaseURL(token string, model string, baseURL string, system string) *OpenRouter {
	return &OpenRouter{
		client:  &http.Client{},
		token:   token,
		model:   model,
		BaseURL: baseURL,
		system:  system,
	}
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (o *OpenRouter) Completion(ctx context.Context, req ports.CompletionRequest, callback func(resp ports.CompletionResponse) error) error {
	reqBody := struct {
		Messages []Message `json:"messages"`
		Model    string    `json:"model,omitempty"`
	}{
		Messages: []Message{
			{
				Role:    "system",
				Content: o.system,
			},
			{
				Role:    "user",
				Content: req.Prompt,
			},
		},
		Model: o.model,
	}

	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	if err := encoder.Encode(reqBody); err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", o.BaseURL+"/completions", buffer)
	if err != nil {
		return err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+o.token)

	httpResp, err := o.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", httpResp.StatusCode)
	}

	respBody := struct {
		Choices []struct {
			Text string `json:"text"`
		} `json:"choices"`
	}{}

	decoder := json.NewDecoder(httpResp.Body)
	if err := decoder.Decode(&respBody); err != nil {
		return err
	}

	if len(respBody.Choices) == 0 {
		return nil
	}

	return callback(ports.CompletionResponse{
		Response: respBody.Choices[0].Text,
		Done:     true,
	})
}
