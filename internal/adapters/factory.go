package adapters

import (
	"errors"

	"github.com/bernardo-bruning/ollama-copilot/internal/ports"
)

var ErrUnknownProvider = errors.New("unknown provider")

func NewProvider(provider string, model string, token string, numPredict int, system string, templateStr string) (ports.Provider, error) {
	switch provider {
	case "ollama":
		return NewOllama(model, numPredict, system)
	case "openrouter":
		return NewOpenRouter(token, model, system, templateStr), nil
	case "deepseek":
		client := NewDefaultHttpClient("https://api.deepseek.com", token)
		return NewDeepSeek(model, numPredict, client), nil
	case "mistral":
		return NewMistral(token, model, system), nil
	case "openai":
		return NewOpenAI(token, model, numPredict), nil
	default:
		return nil, ErrUnknownProvider
	}
}
