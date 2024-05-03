package middleware

import (
	"log"
	"net/http"
	"time"
)

type ResponseWriterLogged struct {
	http.ResponseWriter
	Status int
}

func (w *ResponseWriterLogged) WriteHeader(status int) {
	w.Status = status
	w.ResponseWriter.WriteHeader(status)
}

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		record := ResponseWriterLogged{w, http.StatusOK}
		log.Printf("request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(&record, r)
		log.Printf("response: %s %s %d %s", r.Method, r.URL.Path, record.Status, time.Since(start))
	})
}
