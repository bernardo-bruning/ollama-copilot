package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/bernardo-bruning/ollama-copilot/internal/handlers"
)

func TestTokenHandler_ServeHTTP(t *testing.T) {
	handler := handlers.NewTokenHandler()

	req := httptest.NewRequest(http.MethodGet, "/copilot_internal/v2/token", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var token handlers.TokenResponse
	if err := json.NewDecoder(w.Body).Decode(&token); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	expected := handlers.Token()

	if !reflect.DeepEqual(token, expected) {
		t.Errorf("expected response to be %v, got %v", expected, token)
	}
}
