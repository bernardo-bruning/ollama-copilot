package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/bernardo-bruning/ollama-copilot/internal/ports"
)

type OpenRouter struct {
	client   *http.Client
	token    string
	model    string
	system   string
	BaseURL  string
	template *template.Template
}

func NewOpenRouter(token string, model string, system string, templateStr string) *OpenRouter {
	return NewOpenRouterWithBaseURL(token, model, "https://openrouter.ai/api/v1", system, templateStr)
}

func NewOpenRouterWithBaseURL(token string, model string, baseURL string, system string, templateStr string) *OpenRouter {
	var templ *template.Template
	if templateStr != "" {
		templ, _ = template.New("prompt").Parse(templateStr)
	}
	return &OpenRouter{
		client:   &http.Client{},
		token:    token,
		model:    model,
		BaseURL:  baseURL,
		system:   system,
		template: templ,
	}
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Completion is the completion handler for OpenRouter
func (o *OpenRouter) Completion(ctx context.Context, req ports.CompletionRequest, callback func(resp ports.CompletionResponse) error) error {
	prompt, err := o.generatePrompt(req)
	if err != nil {
		return err
	}

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
				Content: prompt,
			},
		},
		Model: o.model,
	}

	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	if err := encoder.Encode(reqBody); err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", o.BaseURL+"/chat/completions", buffer)
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
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
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
		Response: respBody.Choices[0].Message.Content,
		Done:     true,
	})
}

func (o *OpenRouter) generatePrompt(req ports.CompletionRequest) (string, error) {
	if o.template != nil {
		var buf = new(bytes.Buffer)
		err := o.template.Execute(buf, map[string]string{
			"Prefix":     req.Prompt,
			"Suffix":     req.Suffix,
			"Prompt":     req.Prompt,
			"Prediction": req.Prompt,
		})
		if err != nil {
			return "", err
		}
		return buf.String(), nil
	}
	return req.Prompt, nil
}