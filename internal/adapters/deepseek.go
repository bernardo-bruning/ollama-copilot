package adapters

import (
	"context"

	"github.com/bernardo-bruning/ollama-copilot/internal/ports"
)

type DeepSeekRequest struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	MaxTokens   int     `json:"max_tokens"`
	Stream      bool    `json:"stream"`
	Suffix      string  `json:"suffix"`
	Temperature float64 `json:"temperature"`
	TopP        int     `json:"top_p"`
}
type DeepSeekResponse struct {
	Id      string `json:"id"`
	Choices struct {
		Text string `json:"text"`
	} `json:"choices"`
}

type DeepSeek struct {
	client    HttpClient
	model     string
	maxTokens int
}

// Completion implements [ports.Provider].
func (d *DeepSeek) Completion(ctx context.Context, req ports.CompletionRequest, callback func(resp ports.CompletionResponse) error) error {
	var deepSeekResponse DeepSeekResponse
	deepSeekRequest := DeepSeekRequest{
		Model:       d.model,
		Prompt:      req.Prompt,
		Suffix:      req.Suffix,
		Temperature: req.Temperature,
		TopP:        req.TopP,
		MaxTokens:   d.maxTokens,
	}

	err := d.client.Post("/beta/completions", deepSeekRequest, &deepSeekResponse)
	if err != nil {
		return err
	}

	return callback(ports.CompletionResponse{
		Response: deepSeekResponse.Choices.Text,
		Done:     true,
	})
}

func NewDeepSeek(model string, maxTokens int, httpClient HttpClient) ports.Provider {
	return &DeepSeek{
		client:    httpClient,
		model:     model,
		maxTokens: maxTokens,
	}
}
