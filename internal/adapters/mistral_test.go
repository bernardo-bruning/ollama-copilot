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

func TestMistral(t *testing.T) {
	tests := []struct {
		name             string
		expectedResponse string
		expectedText     string
		expectedPrompt   string
		expectedSuffix   string
		statusCode       int
		expectError      bool
	}{
		{
			name: "Successful completion",
			expectedResponse: `
{
  "choices": [
    {
      "text": "The theory of relativity, developed by Albert Einstein..."
    }
  ]
}`,
			expectedText:   "The theory of relativity, developed by Albert Einstein...",
			expectedPrompt: "Explain the theory of relativity in simple terms.",
			expectedSuffix: "",
			statusCode:     http.StatusOK,
			expectError:    false,
		},
		{
			name: "Successful completion with suffix",
			expectedResponse: `
{
  "choices": [
    {
      "text": "code completion result"
    }
  ]
}`,
			expectedText:   "code completion result",
			expectedPrompt: "func main() {",
			expectedSuffix: "}",
			statusCode:     http.StatusOK,
			expectError:    false,
		},
		{
			name:             "Empty choices",
			expectedResponse: `{"choices": []}`,
			expectedText:     "",
			expectedPrompt:   "Hello",
			expectedSuffix:   "",
			statusCode:       http.StatusOK,
			expectError:      false,
		},
		{
			name:             "Server error",
			expectedResponse: `{"error": "Internal Server Error"}`,
			expectedText:     "",
			expectedPrompt:   "Hello",
			expectedSuffix:   "",
			statusCode:       http.StatusInternalServerError,
			expectError:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/fim/completions" {
					t.Errorf("expected path /fim/completions, got %s", r.URL.Path)
				}
				if r.Method != "POST" {
					t.Errorf("expected method POST, got %s", r.Method)
				}

				var reqBody struct {
					Prompt string `json:"prompt"`
					Suffix string `json:"suffix"`
					Model  string `json:"model"`
				}
				err := json.NewDecoder(r.Body).Decode(&reqBody)
				if err != nil {
					t.Fatalf("failed to decode request body: %v", err)
				}

				if reqBody.Prompt != tt.expectedPrompt {
					t.Errorf("expected prompt '%s', got '%s'", tt.expectedPrompt, reqBody.Prompt)
				}
				if reqBody.Suffix != tt.expectedSuffix {
					t.Errorf("expected suffix '%s', got '%s'", tt.expectedSuffix, reqBody.Suffix)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.expectedResponse))
			}))
			defer server.Close()

			// System prompt is ignored in FIM but passed in constructor
			mistral := adapters.NewMistralWithBaseURL("<token>", "mistral-tiny", server.URL, "ignored-system")

			req := ports.CompletionRequest{
				Prompt: tt.expectedPrompt, // Use Prompt
				Suffix:    tt.expectedSuffix,
			}

			var result string
			err := mistral.Completion(context.Background(), req, func(resp ports.CompletionResponse) error {
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
