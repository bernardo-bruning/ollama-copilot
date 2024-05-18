package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"flag"
	"log"
	"math/big"
	"net/http"
	"text/template"
	"time"

	"github.com/bernardo-bruning/ollama-copilot/internal"
	"github.com/bernardo-bruning/ollama-copilot/internal/handlers"
	"github.com/bernardo-bruning/ollama-copilot/internal/middleware"

	"github.com/ollama/ollama/api"
)

var (
	port        = flag.String("port", ":11437", "Port to listen on")
	portSSL     = flag.String("port-ssl", ":11436", "Port to listen on")
	proxyPort   = flag.String("proxy-port", ":11435", "Proxy port to listen on")
	cert        = flag.String("cert", "", "Certificate file path *.crt")
	key         = flag.String("key", "", "Key file path *.key")
	model       = flag.String("model", "codellama:code", "LLM model to use")
	numPredict  = flag.Int("num-predict", 50, "Number of predictions to return")
	templateStr = flag.String("template", "<PRE> {{.Prefix}} <SUF> {{.Suffix}} <MID>", "Fill-in-middle template to apply in prompt")
)

// main is the entrypoint for the program.
func main() {
	flag.Parse()
	api, err := api.ClientFromEnvironment()

	if err != nil {
		log.Fatalf("error initialize api: %s", err.Error())
		return
	}

	templ, err := template.New("prompt").Parse(*templateStr)
	if err != nil {
		log.Fatalf("error parsing template: %s", err.Error())
		return
	}

	mux := http.NewServeMux()

	mux.Handle("/health", handlers.NewHealthHandler())
	mux.Handle("/copilot_internal/v2/token", handlers.NewTokenHandler())
	mux.Handle("/v1/engines/copilot-codex/completions", handlers.NewCompletionHandler(api, *model, templ, *numPredict))

	handler := middleware.LogMiddleware(mux)
	go internal.Proxy(*proxyPort, *portSSL)

	go listenAndServeTLS(*portSSL, *cert, *key, handler)

	listenAndServe(*port, handler)
}

func listenAndServe(port string, mux http.Handler) {
	err := http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatalf("error listening: %s", err.Error())
	}
}

func listenAndServeTLS(portSSL string, cert string, key string, mux http.Handler) {
	server := http.Server{
		Addr:      portSSL,
		Handler:   mux,
		TLSConfig: &tls.Config{Certificates: []tls.Certificate{}, MinVersion: tls.VersionTLS13, MaxVersion: tls.VersionTLS13},
	}

	if cert == "" || key == "" {
		selfAssignCertificate, err := selfAssignCertificate()
		if err != nil {
			log.Fatalf("error self assigning certificate: %s", err.Error())
		}

		server.TLSConfig.Certificates = append(server.TLSConfig.Certificates, selfAssignCertificate)
	}

	err := server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatalf("error listening: %s", err.Error())
	}
}

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
