package main

import (
	"flag"

	"github.com/bernardo-bruning/ollama-copilot/internal"
)

var (
	port         = flag.String("port", ":11437", "Port to listen on")
	proxyPort    = flag.String("proxy-port", ":11438", "Proxy port to listen on")
	portSSL      = flag.String("port-ssl", ":11436", "Port to listen on")
	proxyPortSSL = flag.String("proxy-port-ssl", ":11435", "Proxy port to listen on")
	cert         = flag.String("cert", "", "Certificate file path *.crt")
	key          = flag.String("key", "", "Key file path *.key")
	model        = flag.String("model", "codellama:code", "LLM model to use")
	numPredict   = flag.Int("num-predict", 50, "Number of predictions to return")
	templateStr  = flag.String("template", "<PRE> {{.Prefix}} <SUF> {{.Suffix}} <MID>", "Fill-in-middle template to apply in prompt")
	system       = flag.String("system", "You are a helpful coding assistant. Respond with autocomplete code only, without explanations or comments.", "The system parameter is use to provide system-level instructions to guide the model's behavior throughout the conversation")
)

// main is the entrypoint for the program.
func main() {
	flag.Parse()
	server := &internal.Server{
		PortSSL:     *portSSL,
		Port:        *port,
		Certificate: *cert,
		Key:         *key,
		Template:    *templateStr,
		Model:       *model,
		NumPredict:  *numPredict,
		System:      *system,
	}

	go internal.Proxy(*proxyPortSSL, *portSSL)
	go internal.Proxy(*proxyPort, *port)

	go server.Serve()
	server.ServeTLS()
}
