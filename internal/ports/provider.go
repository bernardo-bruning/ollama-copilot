package ports

import "context"

type Provider interface {
	Completion(ctx context.Context) error
}

type CompletionRequest struct {
	Prompt      string
	Temperature float64
	TopP        int
	Stop        []string
}

type CompletionResponse struct {
	Response string
	Done     bool
}
