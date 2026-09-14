package adapters

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/bernardo-bruning/ollama-copilot/internal/ports"
	"github.com/ollama/ollama/api"
)

type Ollama struct {
	model      string
	numPredict int
	numCtx     int
	system     string
	client     *api.Client
}

type authTransport struct {
	token string
	base  http.RoundTripper
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.token != "" {
		req.Header.Set("Authorization", "Bearer "+t.token)
	}
	return t.base.RoundTrip(req)
}

// NewOllama creates a new Ollama adapter
func NewOllama(model string, token string, numPredict int, numCtx int, system string) (ports.Provider, error) {
	return newOllamaAPI("OLLAMA_HOST", "http://127.0.0.1:11434", model, token, numPredict, numCtx, system)
}

// NewLlmman creates an adapter for llmman (https://github.com/llmmanorg/llmman),
// a local model runner that serves the Ollama API on port 17434.
func NewLlmman(model string, token string, numPredict int, numCtx int, system string) (ports.Provider, error) {
	return newOllamaAPI("LLMMAN_HOST", "http://127.0.0.1:17434", model, token, numPredict, numCtx, system)
}

// newOllamaAPI creates an adapter for any server speaking the Ollama API,
// reading the host from hostEnv and falling back to defaultHost.
func newOllamaAPI(hostEnv string, defaultHost string, model string, token string, numPredict int, numCtx int, system string) (ports.Provider, error) {
	host := os.Getenv(hostEnv)
	if host == "" {
		host = defaultHost
	}

	if !strings.Contains(host, "://") {
		host = "http://" + host
	}

	baseURL, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{
		Transport: &authTransport{
			token: token,
			base:  http.DefaultTransport,
		},
	}

	client := api.NewClient(baseURL, httpClient)

	return &Ollama{
		model:      model,
		numPredict: numPredict,
		numCtx:     numCtx,
		client:     client,
		system:     system,
	}, nil
}

// Completion is the completion handler for Ollama
func (o *Ollama) Completion(ctx context.Context, req ports.CompletionRequest, callback func(resp ports.CompletionResponse) error) error {
	generate := api.GenerateRequest{
		Model:  o.model,
		Prompt: req.Prompt,
		Suffix: req.Suffix,
		System: o.system,
		Options: map[string]interface{}{
			"temperature": req.Temperature,
			"top_p":       req.TopP,
			"stop":        append(req.Stop, "<EOT>"),
			"num_predict": o.numPredict,
		},
	}

	if o.numCtx > 0 {
		generate.Options["num_ctx"] = o.numCtx
	}

	return o.client.Generate(ctx, &generate, func(resp api.GenerateResponse) error {
		return callback(ports.CompletionResponse{
			Response: resp.Response,
			Done:     resp.Done,
		})
	})
}
