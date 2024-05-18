package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// TokenResponse is the response returned by the TokenHandler.
type TokenResponse struct {
	AnnotationEnabled                  bool     `json:"annotation_enabled"`
	ChatEnabled                        bool     `json:"chat_enabled"`
	CodeQuoteEnabled                   bool     `json:"code_quote_enabled"`
	CopilotIdeAgentChatGpt4SmallPrompt bool     `json:"copilot_ide_agent_chat_gpt4_small_prompt"`
	CopilotIgnoreEnabled               bool     `json:"copilotignore_enabled"`
	ExpiresAt                          int64    `json:"expires_at"`
	IndividualChatEnabled              bool     `json:"individual_chat_enabled"`
	NesEnabled                         bool     `json:"nes_enabled"`
	OrganizationList                   []string `json:"organization_list"`
	Prompt8k                           bool     `json:"prompt_8k"`
	PublicSuggestions                  string   `json:"public_suggestions"`
	RefreshIn                          int64    `json:"refresh_in"`
	Sku                                string   `json:"sku"`
	SnippyLoadTestEnabled              bool     `json:"snippy_load_test_enabled"`
	Telemetry                          string   `json:"telemetry"`
	Token                              string   `json:"token"`
	TrackingId                         string   `json:"tracking_id"`
	VscElectronFetcher                 bool     `json:"vsc_electron_fetcher"`
}

// TokenHandler is an http.Handler that returns a token.
type TokenHandler struct {
}

// NewTokenHandler returns a new TokenHandler.
func NewTokenHandler() *TokenHandler {
	return &TokenHandler{}
}

// ServeHTTP implements http.Handler.
func (t *TokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	token := Token()

	w.Header().Set("content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(token)
	if err != nil {
		log.Printf("error encoding: %s", err.Error())
	}
}

func Token() TokenResponse {
	return TokenResponse{
		AnnotationEnabled:                  false,
		ChatEnabled:                        false,
		CodeQuoteEnabled:                   true,
		CopilotIdeAgentChatGpt4SmallPrompt: false,
		CopilotIgnoreEnabled:               false,
		ExpiresAt:                          time.Now().Unix() + 7200,
		IndividualChatEnabled:              false,
		NesEnabled:                         true,
		OrganizationList:                   []string{},
		Prompt8k:                           true,
		PublicSuggestions:                  "public_suggestions",
		RefreshIn:                          1500,
		Sku:                                "sku",
		SnippyLoadTestEnabled:              true,
		Telemetry:                          "disabled",
		Token:                              "tid=aaaaaaaaaaaaaaaaaaaaaa",
		TrackingId:                         "aaaaaaaaaaaaaaaaaaaaaa",
		VscElectronFetcher:                 true,
	}
}
