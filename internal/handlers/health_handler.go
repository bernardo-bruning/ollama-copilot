package handlers

import (
	"log"
	"net/http"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Ollama copilot is running"))
	if err != nil {
		log.Printf("error writing response: %s", err.Error())
	}
}
