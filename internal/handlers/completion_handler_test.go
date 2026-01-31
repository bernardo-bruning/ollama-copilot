package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bernardo-bruning/ollama-copilot/internal/handlers"
	"github.com/bernardo-bruning/ollama-copilot/internal/ports"
)

type mockProvider struct{}

func (m *mockProvider) Completion(ctx context.Context, req ports.CompletionRequest, callback func(resp ports.CompletionResponse) error) error {
	return nil
}

func TestCompletionHandler_InvalidJSON(t *testing.T) {
	provider := &mockProvider{}
	h := handlers.NewCompletionHandler(provider)

	// Invalid JSON body
	reqBody := strings.NewReader(`{ "invalid": `)
	req, err := http.NewRequest("POST", "/completion", reqBody)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}
