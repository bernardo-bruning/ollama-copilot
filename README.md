# Ollama Copilot

Proxy that allows you to use ollama as a copilot like Github copilot

## Installation

1. Install binaries
```bash
curl https://raw.githubusercontent.com/bernardo-bruning/ollama-copilot/master/scripts/install.sh | sh
```

2. Running
```bash
ollama-copilot
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
