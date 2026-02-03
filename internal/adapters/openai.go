package adapters

import (
	"context"

	"github.com/bernardo-bruning/ollama-copilot/internal/ports"
)

type OpenAI struct {
	token     string
	model     string
	maxTokens int
}

// Completion implements [ports.Provider].
func (o *OpenAI) Completion(ctx context.Context, req ports.CompletionRequest, callback func(resp ports.CompletionResponse) error) error {
	client := NewHttpClient("https://api.openai.com/v1", map[string]string{"Authorization": "Bearer " + o.token})
	openAIRequest := newOpenAIRequest(req)
	openAIResponse := newOpenAIResponse()
	err := client.Post("v1/completions", openAIRequest, &openAIResponse)
	if err != nil {
		return err
	}

	return callback(openAIResponse.ToPorts())
}

func NewOpenAI(token string, model string, maxTokens int) ports.Provider {
	return &OpenAI{token, model, maxTokens}
}

type openAIRequest struct {
	Prompt    string `json:"prompt"`
	Model     string `json:"model"`
	MaxTokens int    `json:"max_tokens"`
	Suffix    string `json:"suffix"`
	// TODO: add temperature and top_p
}

func newOpenAIRequest(req ports.CompletionRequest) *openAIRequest {
	return &openAIRequest{
		Prompt: req.Prompt,
		Suffix: req.Suffix,
	}
}

type openAIResponse struct {
	Id      string           `json:"id"`
	Created int64            `json:"created"`
	Choices []ChoiceResponse `json:"choices"`
}

type ChoiceResponse struct {
	Text         string `json:"text"`
	Index        int    `json:"index"`
	FinishReason string `json:"finish_reason"`
}

func newOpenAIResponse() *openAIResponse {
	return &openAIResponse{}
}

func (o *openAIResponse) ToPorts() ports.CompletionResponse {
	return ports.CompletionResponse{
		Response: o.GetText(),
		Done:     true,
	}
}

func (o *openAIResponse) GetText() string {
	if len(o.Choices) == 0 {
		return ""
	}

	return o.Choices[0].Text
}
