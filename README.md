# Ollama Copilot

Proxy that allows you to use ollama as a copilot like Github copilot

![Video presentation](presentation.gif)

## Installation

### Ollama

Ensure [ollama](https://ollama.com/download/linux) is installed:

```bash
curl -fsSL https://ollama.com/install.sh | sh
```

Or follow the [manual install](https://github.com/ollama/ollama/blob/main/docs/linux.md#manual-install).

#### Models

To use the default model expected by `ollama-copilot`:

```bash
ollama pull codellama:code
```

### ollama-copilot

```bash
go install github.com/bernardo-bruning/ollama-copilot@latest
```

### Running

Ensure your `$PATH` includes `$HOME/go/bin` or `$GOPATH/bin`.
For example, in `~/.bashrc` or `~/.zshrc`:

```bash
export PATH="$HOME/go/bin:$GOPATH/bin:$PATH"
```

```bash
ollama-copilot
```

or if you are hosting ollama in a container or elsewhere
```bash
OLLAMA_HOST="http://192.168.133.7:11434" ollama-copilot
```

## Configure IDE

### Neovim

1. Install [copilot.vim](https://github.com/github/copilot.vim)
1. Configure variables

```vim
let g:copilot_proxy = 'http://localhost:11435'
let g:copilot_proxy_strict_ssl = v:false
```

### VScode

1. Install [copilot extension](https://marketplace.visualstudio.com/items?itemName=GitHub.copilot)
1. Sign-in or sign-up in github
1. Configure open [settings](https://code.visualstudio.com/docs/getstarted/settings) config and insert

```json
{
    "github.copilot.advanced": {
        "debug.overrideProxyUrl": "http://localhost:11437"
    },
    "http.proxy": "http://localhost:11435",
    "http.proxyStrictSSL": false
}
```

### Emacs

(experimental)

1. Install [copilot-emacs](https://github.com/copilot-emacs/copilot.el)
1. Configure the proxy

```elisp
(use-package copilot
  :straight (:host github :repo "copilot-emacs/copilot.el" :files ("*.el"))  ;; if you don't use "straight", install otherwise
  :ensure t
  ;; :hook (prog-mode . copilot-mode)
  :bind (
         ("C-<tab>" . copilot-accept-completion)
         )
  :config
  (setq copilot-network-proxy '(:host "127.0.0.1" :port 11434 :rejectUnauthorized :json-false))
  )
```


## Roadmap

- [x] Enable completions APIs usage; fill in the middle.
- [x] Enable flexible configuration model (Currently only supported llamacode:code).
- [x] Create self-installing functionality.
- [ ] Windows setup
- [ ] Documentation on how to use.
