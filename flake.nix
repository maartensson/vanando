{
  description = "vanando image webservice";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.11";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {self, nixpkgs, flake-utils, ...}:
  flake-utils.lib.eachDefaultSystem (system: {
    packages.default = nixpkgs.legacyPackages.${system}.buildGoModule {
      pname = "vanando";
      version = "0.0.1";
      src = ./.;
      vendorHash = null;
    };

    apps.default = {
      type = "app";
      program = "${self.packages.${system}.default}/bin/test";
    };
  }) // {
    nixosModules.default = {config, lib, pkgs, ...}: {
      options.services.vanando = {
        enable = lib.mkEnableOption "Enable vanando image webservice";
        port = lib.mkOption {
          type = lib.types.port;
          default = 8080;
          description = "Port to listen on";
        };
      };

      config = lib.mkIf config.services.vanando.enable {
        systemd.services.vanando = {
          description = "Vanando image webservice";
          wantedBy = ["multi-user.target"];
          after = ["network.target"];
          serviceConfig = {
            ExecStart = "${self.packages.${pkgs.system}.default}/bin/test";
            Restart = "always";
            Type = "simple";
            DynamicUser = "yes";
            Environment = [
              "PORT=${toString config.services.vanando.port}"
            ];
          };
        };
      };
    };
  };
}

