package adapters

import (
	"context"

	"github.com/bernardo-bruning/ollama-copilot/internal/ports"
)

type OpenAI struct {
	model     string
	maxTokens int
	client    HttpClient
}

// Completion implements [ports.Provider].
func (o *OpenAI) Completion(ctx context.Context, req ports.CompletionRequest, callback func(resp ports.CompletionResponse) error) error {
	openAIRequest := NewOpenAIRequest(req)
	openAIResponse := NewOpenAIResponse()
	err := o.client.Post("v1/completions", openAIRequest, openAIResponse)
	if err != nil {
		return err
	}

	return callback(openAIResponse.ToPorts())
}

func NewOpenAI(token string, model string, maxTokens int) ports.Provider {
	client := NewDefaultHttpClient("https://api.openai.com", token)
	return &OpenAI{model, maxTokens, client}
}

func NewOpenAIWithClient(model string, maxTokens int, client HttpClient) ports.Provider {
	return &OpenAI{model, maxTokens, client}
}

type OpenAIRequest struct {
	Prompt    string `json:"prompt"`
	Model     string `json:"model"`
	MaxTokens int    `json:"max_tokens"`
	Suffix    string `json:"suffix"`
	// TODO #47:30min add temperature and top_p
}

func NewOpenAIRequest(req ports.CompletionRequest) *OpenAIRequest {
	return &OpenAIRequest{
		Prompt: req.Prompt,
		Suffix: req.Suffix,
	}
}

type OpenAIResponse struct {
	Id      string           `json:"id"`
	Created int64            `json:"created"`
	Choices []ChoiceResponse `json:"choices"`
}

type ChoiceResponse struct {
	Text         string `json:"text"`
	Index        int    `json:"index"`
	FinishReason string `json:"finish_reason"`
}

func NewOpenAIResponse() *OpenAIResponse {
	return &OpenAIResponse{}
}

func (o *OpenAIResponse) ToPorts() ports.CompletionResponse {
	return ports.CompletionResponse{
		Response: o.GetText(),
		Done:     true,
	}
}

func (o *OpenAIResponse) GetText() string {
	if len(o.Choices) == 0 {
		return ""
	}

	return o.Choices[0].Text
}
