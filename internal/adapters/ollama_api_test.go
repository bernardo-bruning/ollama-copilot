package adapters

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bernardo-bruning/ollama-copilot/internal/ports"
)

func TestOllamaCompletionAddsEOT(t *testing.T) {
	// Start a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/generate" {
			t.Errorf("Expected path /api/generate, got %s", r.URL.Path)
		}

		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		options, ok := reqBody["options"].(map[string]interface{})
		if !ok {
			t.Fatalf("Expected options map in request body")
		}

		stop, ok := options["stop"].([]interface{})
		if !ok {
			t.Fatalf("Expected stop slice in options")
		}

		foundEOT := false
		for _, s := range stop {
			if str, ok := s.(string); ok && str == "<EOT>" {
				foundEOT = true
				break
			}
		}

		if !foundEOT {
			t.Errorf("Expected <EOT> in stop sequences, got %v", stop)
		}

		// Respond with something valid so the client doesn't error out immediately
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"response": "ok", "done": true}`))
	}))
	defer server.Close()

	// Set OLLAMA_HOST to the mock server URL
	// The ollama api client might need a clean URL
	t.Setenv("OLLAMA_HOST", server.URL)

	// Create the Ollama adapter
	// Note: NewOllama calls api.ClientFromEnvironment() which reads OLLAMA_HOST
	adapter, err := NewOllama("test-model", 100, "test-system", "{{.Prompt}}")
	if err != nil {
		t.Fatalf("NewOllama failed: %v", err)
	}

	// Call Completion
	req := ports.CompletionRequest{
		Prompt:      "test prompt",
		Stop:        []string{"stop1"},
		Temperature: 0.5,
	}

	err = adapter.Completion(context.Background(), req, func(resp ports.CompletionResponse) error {
		return nil
	})

	if err != nil {
		t.Fatalf("Completion failed: %v", err)
	}
}
