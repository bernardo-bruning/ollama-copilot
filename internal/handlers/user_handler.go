package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type AuthenticationToken struct {
	UserID       string `json:"user_id"`
	Login        string `json:"login"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (t *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	userResponse := GetUserResponse()

	w.Header().Set("content-Type", "application/json")
	w.Header().Set("X-GitHub-Request-Id", "req-id")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(userResponse)
	if err != nil {
		log.Printf("error encoding: %s", err.Error())
	}
}

func GetUserResponse() AuthenticationToken {
	return AuthenticationToken{
		UserID:       "12345678",
		Login:        "octocat",
		AccessToken:  "gho_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		RefreshToken: "ghr_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		Scope:        "copilot read:user",
		ExpiresIn:    3600,
		TokenType:    "Bearer",
	}
}
