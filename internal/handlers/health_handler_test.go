package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bernardo-bruning/ollama-copilot/internal/handlers"
)

func TestHealthHandler(t *testing.T) {
	h := handlers.NewHealthHandler()
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.ServeHTTP)
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := "Ollama copilot is running"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
