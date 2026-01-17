package adapters

import (
	"context"

	"github.com/bernardo-bruning/ollama-copilot/internal/ports"
	"github.com/ollama/ollama/api"
)

type Ollama struct {
	model      string
	numPredict int
	client     *api.Client
}

// NewOllama creates a new Ollama adapter
func NewOllama(client *api.Client, model string, numPredict int) ports.Provider {
	return &Ollama{
		model:      model,
		numPredict: numPredict,
		client:     client,
	}
}

// Completion is the completion handler for Ollama
func (o *Ollama) Completion(ctx context.Context, req ports.CompletionRequest, callback func(resp ports.CompletionResponse) error) error {
	generate := api.GenerateRequest{
		Model:  o.model,
		Prompt: req.Prompt,
		Options: map[string]interface{}{
			"temperature": req.Temperature,
			"top_p":       req.TopP,
			"stop":        req.Stop,
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
