package main

import (
	"log"
	"net/http"
	"ollama-copilot/internal"
	"ollama-copilot/internal/handlers"
	"ollama-copilot/internal/middleware"

	"github.com/ollama/ollama/api"
)

// main is the entrypoint for the program.
func main() {
	api, err := api.ClientFromEnvironment()

	if err != nil {
		log.Fatalf("error initialize api: %s", err.Error())
		return
	}

	mux := http.NewServeMux()

	mux.Handle("/health", handlers.NewHealthHandler())
	mux.Handle("/copilot_internal/v2/token", handlers.NewTokenHandler())
	mux.Handle("/v1/engines/copilot-codex/completions", handlers.NewCompletionHandler(api))

	go internal.Proxy(":8080")

	http.ListenAndServeTLS(":9090", "server.crt", "server.key", middleware.LogMiddleware(mux))
}
