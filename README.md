# Ollama Copilot

Proxy that allows you to use ollama as a copilot like Github copilot

![Video presentation](presentation.gif)

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

1. Install [copilot.vim](https://github.com/github/copilot.vim)
2. Sign-in or sign-up in github
3. Configure variables
```
let g:copilot_proxy = 'http://localhost:11435'
let g:copilot_proxy_strict_ssl = v:false
```

### VScode
1. Install [copilot extension](https://marketplace.visualstudio.com/items?itemName=GitHub.copilot)
2. Configure open [settings](https://code.visualstudio.com/docs/getstarted/settings) config and insert
```
{
    "github.copilot.advanced": {
        "debug.overrideProxyUrl": "http://localhost:11437",
    },
}
```

## Roadmap

- [x] Enable completions APIs usage; fill in the middle.
- [x] Enable flexible configuration model (Currently only supported llamacode:code).
- [x] Create self-installing functionality.
- [ ] Windows setup
- [ ] Documentation on how to use.
