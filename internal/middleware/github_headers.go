package middleware

import (
	"net/http"
)

func GithubHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// The copilot-language-server npm module uses this
		// header to check if the connection to GitHub is working,
		// and otherwise errors out.
		w.Header().Set("x-github-request-id", "foobar")
		next.ServeHTTP(w, r)
	})
}
