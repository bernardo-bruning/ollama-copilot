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
	tests := []struct {
		name             string
		expectedResponse string
		expectedText     string
		expectedPrompt   string
		expectedSystem   string
		statusCode       int
		expectError      bool
	}{
		{
			name: "Successful completion",
			expectedResponse: `
{
  "choices": [
    {
      "text": "The theory of relativity, developed by Albert Einstein, explains how space and time are linked for objects moving at a consistent speed in a straight line. It shows that time can slow down or speed up depending on how fast you move relative to something else."
    }
  ]
}`,
			expectedText:   "The theory of relativity, developed by Albert Einstein, explains how space and time are linked for objects moving at a consistent speed in a straight line. It shows that time can slow down or speed up depending on how fast you move relative to something else.",
			expectedPrompt: "Explain the theory of relativity in simple terms.",
			expectedSystem: "You are a helpful assistant.",
			statusCode:     http.StatusOK,
			expectError:    false,
		},
		{
			name:             "Empty choices",
			expectedResponse: `{"choices": []}`,
			expectedText:     "",
			expectedPrompt:   "Hello",
			expectedSystem:   "",
			statusCode:       http.StatusOK,
			expectError:      false,
		},
		{
			name:             "Server error",
			expectedResponse: `{"error": "Internal Server Error"}`,
			expectedText:     "",
			expectedPrompt:   "Hello",
			expectedSystem:   "",
			statusCode:       http.StatusInternalServerError,
			expectError:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/completions" {
					t.Errorf("expected path /completions, got %s", r.URL.Path)
				}
				if r.Method != "POST" {
					t.Errorf("expected method POST, got %s", r.Method)
				}

				var reqBody struct {
					Messages []struct {
						Role    string `json:"role"`
						Content string `json:"content"`
					}
					Model string `json:"model"`
				}
				err := json.NewDecoder(r.Body).Decode(&reqBody)
				if err != nil {
					t.Fatalf("failed to decode request body: %v", err)
				}

				foundPrompt := false
				for _, msg := range reqBody.Messages {
					if msg.Role == "user" && msg.Content == tt.expectedPrompt {
						foundPrompt = true
					}
					if msg.Role == "system" && msg.Content != tt.expectedSystem {
						t.Errorf("expected system content '%s', got '%s'", tt.expectedSystem, msg.Content)
					}
				}

				if !foundPrompt {
					t.Errorf("expected prompt '%s' not found in messages", tt.expectedPrompt)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.expectedResponse))
			}))
			defer server.Close()

			openRouter := adapters.NewOpenRouterWithBaseURL("<token>", "openrouter-gpt-3.5-turbo", server.URL, tt.expectedSystem)

			req := ports.CompletionRequest{
				Prompt: tt.expectedPrompt,
			}

			var result string
			err := openRouter.Completion(context.Background(), req, func(resp ports.CompletionResponse) error {
				result = resp.Response
				return nil
			})

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if result != tt.expectedText {
				t.Errorf("expected response '%s', got '%s'", tt.expectedText, result)
			}
		})
	}
}