package adapters_test

import (
	"reflect"
	"testing"

	"github.com/bernardo-bruning/ollama-copilot/internal/adapters"
)

func TestFactory(t *testing.T) {
	tests := []struct {
		provider   string
		model      string
		numPredict int
		typeName   string
		err        error
	}{
		{"ollama", "llama2", 128, "Ollama", nil},
		{"openrouter", "openrouter-gpt-3.5-turbo", 256, "OpenRouter", nil},
		{"deepseek", "deepseek-coder", 256, "DeepSeek", nil},
		{"mistral", "mistral-tiny", 256, "Mistral", nil},
		{"unknown", "", 256, "", adapters.ErrUnknownProvider},
	}

	for _, test := range tests {
		provider, err := adapters.NewProvider(test.provider, test.model, "", test.numPredict, "system", "")
		if err != test.err {
			t.Fatalf("expected error %v, got %v", test.err, err)
		}

		if provider == nil && test.err == nil {
			t.Fatalf("expected non-nil provider, got nil")
		}

		if err == nil {
			typeName := reflect.TypeOf(provider).Elem().Name()
			if typeName != test.typeName {
				t.Errorf("expected type %s, got %s", test.typeName, typeName)
			}
		}
	}
}
