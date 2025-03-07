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
            ollama-copilot.model = lib.mkOption {
              type = lib.types.str;
              default = "";
              description = "The model to use (default: don't set the option, let the app choose its default)";
            };
          };

          config = lib.mkIf cfg.enable {
            systemd.services.ollama-copilot =
              let
                model_opt = if cfg.model != "" then " --model " + cfg.model else "";
              in
              {
                description = "Ollama Copilot";
                wantedBy = [ "multi-user.target" ];
                serviceConfig = {
                  Type = "simple";
                  ExecStart = "${lib.getExe self.packages.x86_64-linux.default} ${model_opt}";
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
          vendorHash = "sha256-Lo7IurCQwkQpZe/UncBPn4bz4eASEHJ67RJaOka2Rc4=";
          meta.mainProgram = "ollama-copilot";
        };
    };
}
