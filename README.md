[![PDD status](https://www.0pdd.com/svg?name=bernardo-bruning/ollama-copilot)](https://www.0pdd.com/p?name=bernardo-bruning/ollama-copilot)
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

### DeepSeek

To use DeepSeek:

```bash
ollama-copilot -provider deepseek -token YOUR_DEEPSEEK_API_KEY -model deepseek-coder
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

## Configuration

You can configure the server using command-line flags:

| Flag | Default | Description |
|------|---------|-------------|
| `-port` | `:11437` | HTTP port to listen on |
| `-proxy-port` | `:11438` | HTTP proxy port |
| `-port-ssl` | `:11436` | HTTPS port to listen on |
| `-proxy-port-ssl` | `:11435` | HTTPS proxy port |
| `-cert` | | Certificate file path (`*.crt`) for custom TLS |
| `-key` | | Key file path (`*.key`) for custom TLS |
| `-provider` | `ollama` | Provider to run LLM |
| `-token` | `TOKEN` | Token to pass for provider |
| `-model` | `codellama:code` | LLM model to use |
| `-num-predict` | `50` | Number of tokens to predict |
| `-template` | `<PRE> {{.Prefix}} <SUF> {{.Suffix}} <MID>` | Prompt template for fill-in-middle |
| `-system` | `You are a helpful...` | System prompt to guide the model |

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

### Zed

1. [Open settings](https://zed.dev/docs/configuring-zed) (ctrl + ,)
1. Set up [edit completion proxying](https://github.com/zed-industries/zed/pull/24364):

```json
{
    "features": {
        "edit_prediction_provider": "copilot"
    },
    "show_completions_on_input": true,
    "edit_predictions": {
        "copilot": {
            "proxy": "http://localhost:11435",
            "proxy_no_verify": true
        }
    }
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
