package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/bernardo-bruning/ollama-copilot/internal/ports"
)

type OpenRouter struct {
	client *http.Client
	token  string
}

func NewOpenRouter(token string) *OpenRouter {
	return &OpenRouter{
		client: &http.Client{},
		token:  token,
	}
}

func (o *OpenRouter) Completion(ctx context.Context, req ports.CompletionRequest, callback func(resp ports.CompletionResponse) error) error {
	reqBody := struct {
		Prompt string `json:"prompt"`
	}{
		Prompt: req.Prompt,
	}

	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	if err := encoder.Encode(reqBody); err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://api.openrouter.ai/v1/chat/completions", buffer)
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

	panic("not implemented")
}
