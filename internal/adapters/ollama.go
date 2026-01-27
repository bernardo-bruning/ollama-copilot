package adapters

import (
	"bytes"
	"context"
	"text/template"

	"github.com/bernardo-bruning/ollama-copilot/internal/ports"
	"github.com/ollama/ollama/api"
)

type Ollama struct {
	model      string
	numPredict int
	system     string
	client     *api.Client
	template   *template.Template
}

// NewOllama creates a new Ollama adapter
func NewOllama(model string, numPredict int, system string, templateStr string) (ports.Provider, error) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return nil, err
	}

	templ, err := template.New("prompt").Parse(templateStr)
	if err != nil {
		return nil, err
	}

	return &Ollama{
		model:      model,
		numPredict: numPredict,
		client:     client,
		system:     system,
		template:   templ,
	}, nil
}

// Completion is the completion handler for Ollama
func (o *Ollama) Completion(ctx context.Context, req ports.CompletionRequest, callback func(resp ports.CompletionResponse) error) error {
	prompt, err := o.generatePrompt(req)
	if err != nil {
		return err
	}

	generate := api.GenerateRequest{
		Model:  o.model,
		Prompt: prompt,
		System: o.system,
		Options: map[string]interface{}{
			"temperature": req.Temperature,
			"top_p":       req.TopP,
			"stop":        append(req.Stop, "<EOT>"),
			"num_predict": o.numPredict,
		},
	}

	return o.client.Generate(ctx, &generate, func(resp api.GenerateResponse) error {
		return callback(ports.CompletionResponse{
			Response: resp.Response,
			Done:     resp.Done,
		})
	})
}

func (o *Ollama) generatePrompt(req ports.CompletionRequest) (string, error) {
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