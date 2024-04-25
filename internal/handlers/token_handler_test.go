package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

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

	expected := handlers.TokenResponse{
		AnnotationEnabled:                  false,
		ChatEnabled:                        false,
		CodeQuoteEnabled:                   true,
		CopilotIdeAgentChatGpt4SmallPrompt: false,
		CopilotIgnoreEnabled:               false,
		ExpiresAt:                          time.Now().Unix() + 3600,
		IndividualChatEnabled:              true,
		NesEnabled:                         true,
		OrganizationList:                   []string{"org1", "org2"},
		Prompt8k:                           true,
		PublicSuggestions:                  "public_suggestions",
		RefreshIn:                          time.Now().Unix() + 1800,
		Sku:                                "sku",
		SnippyLoadTestEnabled:              true,
		Telemetry:                          "telemetry",
		Token:                              "token",
		TrackingId:                         "tracking_id",
		VscElectronFetcher:                 true,
	}

	if !reflect.DeepEqual(token, expected) {
		t.Errorf("expected response to be %v, got %v", expected, token)
	}
}
