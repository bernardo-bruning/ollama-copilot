# !/bin/sh
# Install ollama-copilot and its dependencies

# Install binaries
go install github.com/bernardo-bruning/ollama-copilot@latest

# Install ollama
curl -fsSL https://ollama.com/install.sh | sh

# Pull models
ollama pull codellama:code

# Create etc file
sudo mkdir -p /etc/ollama-copilot
sudo chmod o=rwx /etc/ollama-copilot

# Create certificates for copilot
openssl ecparam -genkey -name secp384r1 -out /etc/ollama-copilot/server.key
openssl req -new -x509 -sha256 -key /etc/ollama-copilot/server.key -out /etc/ollama-copilot/server.crt -days 3650 -subj '/CN=ollama-copilot/O=Ollama copilot./C=US'

# Install user folder
sudo mv $HOME/go/bin/ollama-copilot /usr/local/bin/ollama-copilot
