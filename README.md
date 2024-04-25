# Ollama Copilot

Proxy that allows you to use ollama as a copilot like Github copilot

## Installation

1. Install binaries
```bash
go install github.com/bernardo-bruning/ollama-copilot
```

2. Install ollama
```bash
curl -fsSL https://ollama.com/install.sh | sh
```

3. Pull models
```bash
ollama pull codellama:code
````

4. Create certificates for copilot
```bash
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
```

5. Run ollama copilot
```
~/go/bin/ollama-copilot
```

## Configure IDE

### Neovim

1. Install copilot.vim
2. Configure variables
```
let g:copilot_proxy = 'http://localhost:8080'
let g:copilot_proxy_strict_ssl = v:false
```

## Roadmap

- [x] Enable completions APIs usage; fill in the middle.
- [x] Enable flexible configuration model (Currently only supported llamacode:code).
- [ ] Create self-installing functionality.
- [ ] Documentation on how to use.
