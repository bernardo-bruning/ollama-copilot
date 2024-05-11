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

# Install user folder
sudo mv $HOME/go/bin/ollama-copilot /usr/local/bin/ollama-copilot
