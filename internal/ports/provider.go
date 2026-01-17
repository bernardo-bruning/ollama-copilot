package ports

import (
	"context"
)

// Provider is the interface for a completion provider
type Provider interface {
	Completion(ctx context.Context, req CompletionRequest, callback func(resp CompletionResponse) error) error
}

// CompletionRequest is the request for a completion provider
type CompletionRequest struct {
	Prompt      string
	Temperature float64
	TopP        int
	Stop        []string
}

// CompletionResponse is the response for a completion provider
type CompletionResponse struct {
	Response string
	Done     bool
}
