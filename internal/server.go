package internal

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"log"
	"math/big"
	"net/http"
	"text/template"
	"time"

	"github.com/bernardo-bruning/ollama-copilot/internal/handlers"
	"github.com/bernardo-bruning/ollama-copilot/internal/middleware"
	"github.com/ollama/ollama/api"
)

// Server is the main server struct.
type Server struct {
	PortSSL     string
	Port        string
	Certificate string
	Key         string
	Template    string
	Model       string
	NumPredict  int
}

// Serve starts the server.
func (s *Server) Serve() {
	err := http.ListenAndServe(s.Port, s.mux())
	if err != nil {
		log.Fatalf("error listening: %s", err.Error())
	}
}

// ServeTLS starts the server with TLS.
func (s *Server) ServeTLS() {
	server := http.Server{
		Addr:      s.PortSSL,
		Handler:   s.mux(),
		TLSConfig: &tls.Config{Certificates: []tls.Certificate{}, MinVersion: tls.VersionTLS13, MaxVersion: tls.VersionTLS13},
	}

	if s.Certificate == "" || s.Key == "" {
		selfAssignCertificate, err := selfAssignCertificate()
		if err != nil {
			log.Fatalf("error self assigning certificate: %s", err.Error())
		}

		server.TLSConfig.Certificates = append(server.TLSConfig.Certificates, selfAssignCertificate)
	}

	err := server.ListenAndServeTLS(s.Certificate, s.Key)
	if err != nil {
		log.Fatalf("error listening: %s", err.Error())
	}
}

// selfAssignCertificate generates a self-signed certificate for localhost.
func selfAssignCertificate() (tls.Certificate, error) {
	private, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return tls.Certificate{}, err
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "localhost",
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(30, 0, 0),
		KeyUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
		},
		BasicConstraintsValid: true,
	}

	cert, err := x509.CreateCertificate(rand.Reader, template, template, private.Public(), private)

	return tls.Certificate{
		Certificate: [][]byte{cert},
		PrivateKey:  private,
	}, err
}

// mux returns the main mux for the server.
func (s *Server) mux() http.Handler {
	api, err := api.ClientFromEnvironment()

	if err != nil {
		log.Fatalf("error initialize api: %s", err.Error())
		return nil
	}

	templ, err := template.New("prompt").Parse(s.Template)
	if err != nil {
		log.Fatalf("error parsing template: %s", err.Error())
		return nil
	}

	mux := http.NewServeMux()

	mux.Handle("/health", handlers.NewHealthHandler())
	mux.Handle("/copilot_internal/v2/token", handlers.NewTokenHandler())
	mux.Handle("/v1/engines/copilot-codex/completions", handlers.NewCompletionHandler(api, s.Model, templ, s.NumPredict))

	return middleware.LogMiddleware(mux)
}
