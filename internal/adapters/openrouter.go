package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/bernardo-bruning/ollama-copilot/internal/ports"
)

type OpenRouter struct {
	client  *http.Client
	token   string
	model   string
	BaseURL string
}

func NewOpenRouter(token string, model string) *OpenRouter {
	return NewOpenRouterWithBaseURL(token, model, "https://openrouter.ai/api/v1")
}

func NewOpenRouterWithBaseURL(token string, model string, baseURL string) *OpenRouter {
	return &OpenRouter{
		client:  &http.Client{},
		token:   token,
		model:   model,
		BaseURL: baseURL,
	}
}

func (o *OpenRouter) Completion(ctx context.Context, req ports.CompletionRequest, callback func(resp ports.CompletionResponse) error) error {
	reqBody := struct {
		Prompt string `json:"prompt"`
		Model  string `json:"model,omitempty"`
	}{
		Prompt: req.Prompt,
		Model:  o.model,
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
