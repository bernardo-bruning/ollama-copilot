package adapters

import (
	"context"

	"github.com/bernardo-bruning/ollama-copilot/internal/ports"
)

type Mistral struct {
	client HttpClient
	model  string
	system string
}

func NewMistral(token string, model string, system string) *Mistral {
	return NewMistralWithBaseURL(token, model, "https://api.mistral.ai/v1", system)
}

func NewMistralWithBaseURL(token string, model string, baseURL string, system string) *Mistral {
	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}
	return &Mistral{
		client: NewHttpClient(baseURL, headers),
		model:  model,
		system: system,
	}
}

func (m *Mistral) Completion(ctx context.Context, req ports.CompletionRequest, callback func(resp ports.CompletionResponse) error) error {
	reqBody := struct {
		Model       string   `json:"model"`
		Prompt      string   `json:"prompt"`
		Suffix      string   `json:"suffix,omitempty"`
		Temperature float64  `json:"temperature,omitempty"`
		TopP        int      `json:"top_p,omitempty"`
		Stop        []string `json:"stop,omitempty"`
		Stream      bool     `json:"stream"`
	}{
		Model:       m.model,
		Prompt:      req.Prompt,
		Suffix:      req.Suffix,
		Temperature: req.Temperature,
		TopP:        req.TopP,
		Stop:        req.Stop,
		Stream:      false,
	}

	respBody := struct {
		Choices []struct {
			Text string `json:"text"`
		} `json:"choices"`
	}{}

	err := m.client.Post("/fim/completions", reqBody, &respBody)
	if err != nil {
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
