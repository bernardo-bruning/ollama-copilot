{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
  };

  outputs =
    { self, nixpkgs, ... }:
    {
      nixosModules.default =
        { config, lib, ... }:
        let
          cfg = config.ollama-copilot;
        in
        {
          options = {
            ollama-copilot.enable = lib.mkEnableOption "Enable ollama-copilot";
            ollama-copilot.cert = lib.mkOption {
              type = lib.types.str;
              default = "";
              description = "Certificate file path *.crt";
            };
            ollama-copilot.key = lib.mkOption {
              type = lib.types.str;
              default = "";
              description = "Key file path *.key";
            };
            ollama-copilot.model = lib.mkOption {
              type = lib.types.str;
              default = "";
              description = "LLM model to use (default: codellama:code)";
            };
            ollama-copilot.num-predict = lib.mkOption {
              type = lib.types.int;
              default = 0;
              description = "Number of predictions to return (default: 50)";
            };
            ollama-copilot.port = lib.mkOption {
              type = lib.types.str;
              default = "";
              description = "Port to listen on (default: :11437)";
            };
            ollama-copilot.port-ssl = lib.mkOption {
              type = lib.types.str;
              default = "";
              description = "Port to listen on (default: :11436)";
            };
            ollama-copilot.proxy-port = lib.mkOption {
              type = lib.types.str;
              default = "";
              description = "Proxy port to listen on (default: :11438)";
            };
            ollama-copilot.proxy-port-ssl = lib.mkOption {
              type = lib.types.str;
              default = "";
              description = "Proxy port to listen on (default: :11435)";
            };
            ollama-copilot.system = lib.mkOption {
              type = lib.types.str;
              default = "";
              description = "System-level instructions to guide the model's behavior";
            };
            ollama-copilot.template = lib.mkOption {
              type = lib.types.str;
              default = "";
              description = "Fill-in-middle template to apply in prompt";
            };
          };

          config = lib.mkIf cfg.enable {
            systemd.services.ollama-copilot =
              let
                cert_opt = if cfg.cert != "" then " -cert " + cfg.cert else "";
                key_opt = if cfg.key != "" then " -key " + cfg.key else "";
                model_opt = if cfg.model != "" then " -model " + cfg.model else "";
                num_predict_opt = if cfg.num-predict != 0 then " -num-predict " + toString cfg.num-predict else "";
                port_opt = if cfg.port != "" then " -port " + cfg.port else "";
                port_ssl_opt = if cfg.port-ssl != "" then " -port-ssl " + cfg.port-ssl else "";
                proxy_port_opt = if cfg.proxy-port != "" then " -proxy-port " + cfg.proxy-port else "";
                proxy_port_ssl_opt = if cfg.proxy-port-ssl != "" then " -proxy-port-ssl " + cfg.proxy-port-ssl else "";
                system_opt = if cfg.system != "" then " -system " + cfg.system else "";
                template_opt = if cfg.template != "" then " -template " + cfg.template else "";
              in
              {
                description = "Ollama Copilot";
                wantedBy = [ "multi-user.target" ];
                serviceConfig = {
                  Type = "simple";
                  ExecStart = "${lib.getExe self.packages.x86_64-linux.default}${cert_opt}${key_opt}${model_opt}${num_predict_opt}${port_opt}${port_ssl_opt}${proxy_port_opt}${proxy_port_ssl_opt}${system_opt}${template_opt}";
                  Restart = "always";
                };
              };
          };
        };

      packages.x86_64-linux.default =
        let
          pkgs = import nixpkgs { system = "x86_64-linux"; };
        in
        pkgs.buildGoModule rec {
          name = "ollama-copilot";
          src = ./.;
          vendorHash = "sha256-g27MqS3qk67sve/jexd07zZVLR+aZOslXrXKjk9BWtk=";
          meta.mainProgram = "ollama-copilot";
        };
    };
}
