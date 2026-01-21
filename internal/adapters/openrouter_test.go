package adapters_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bernardo-bruning/ollama-copilot/internal/adapters"
	"github.com/bernardo-bruning/ollama-copilot/internal/ports"
)

func TestOpenRouter(t *testing.T) {
	expectedResponse := `
{
  "id": "cmpl-6aX9b2fJ7Qz8YvN3LpR4TqW1eZxYvB",
  "object": "text_completion",
  "created": 1685600000,
  "model": "openrouter-gpt-3.5-turbo",
  "choices": [
    {
      "text": "The theory of relativity, developed by Albert Einstein, explains how space and time are linked for objects moving at a consistent speed in a straight line. It shows that time can slow down or speed up depending on how fast you move relative to something else.",
      "index": 0,
      "logprobs": {
        "tokens": [
          "The"
        ],
        "token_logprobs": [
          -0.01
        ],
        "top_logprobs": [
          {}
        ],
        "text_offset": [
          0
        ]
      },
      "finish_reason": "stop",
      "native_finish_reason": "stop",
      "reasoning": "The completion ends naturally after explaining the concept."
    }
  ],
  "provider": "openrouter",
  "system_fingerprint": "abc123def456ghi789",
  "usage": {
    "prompt_tokens": 7,
    "completion_tokens": 54,
    "total_tokens": 61
  }
}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/completions" {
			t.Errorf("expected path /completions, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("expected method POST, got %s", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer <token>" {
			t.Errorf("expected Authorization header Bearer <token>, got %s", r.Header.Get("Authorization"))
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}

		var reqBody map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if reqBody["prompt"] != "Explain the theory of relativity in simple terms." {
			t.Errorf("expected prompt 'Explain the theory of relativity in simple terms.', got %v", reqBody["prompt"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(expectedResponse))
	}))
	defer server.Close()

	openRouter := adapters.NewOpenRouterWithBaseURL("<token>", "openrouter-gpt-3.5-turbo", server.URL)

	req := ports.CompletionRequest{
		Prompt: "Explain the theory of relativity in simple terms.",
	}

	var result string
	err := openRouter.Completion(context.Background(), req, func(resp ports.CompletionResponse) error {
		result = resp.Response
		return nil
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedText := "The theory of relativity, developed by Albert Einstein, explains how space and time are linked for objects moving at a consistent speed in a straight line. It shows that time can slow down or speed up depending on how fast you move relative to something else."
	if result != expectedText {
		t.Errorf("expected response '%s', got '%s'", expectedText, result)
	}
}
