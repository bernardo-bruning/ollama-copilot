package adapters

import (
	"errors"

	"github.com/bernardo-bruning/ollama-copilot/internal/ports"
)

var ErrUnknownProvider = errors.New("unknown provider")

func NewProvider(provider string, model string, token string, numPredict int) (ports.Provider, error) {
	switch provider {
	case "ollama":
		return NewOllama(model, numPredict)
	case "openrouter":
		return NewOpenRouter(token, model), nil
	default:
		return nil, ErrUnknownProvider
	}
}
