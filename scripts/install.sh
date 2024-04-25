# !/bin/sh
# Install ollama-copilot and its dependencies

# Install binaries
go install github.com/bernardo-bruning/ollama-copilot

# Install ollama
curl -fsSL https://ollama.com/install.sh | sh

# Pull models
ollama pull codellama:code

# Create certificates for copilot
openssl ecparam -genkey -name secp384r1 -out server.key
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650 -subj '/CN=ollama-copilot/O=Ollama copilot./C=US'
