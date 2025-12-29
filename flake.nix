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
      vendorHash = "sha256-XxxoAoHKQLChs09BE05AEztzzhuQhPSKu7mCD3zsu34=";
    };

    apps.default = {
      type = "app";
      program = "${self.packages.${system}.default}/bin/test";
    };
  }) // {
    nixosModules.default = {config, lib, pkgs, ...}: let 
      system = pkgs.stdenv.hostPlatform.system;
      vanando = self.packages.${system}.default;
      cfg = config.services.vanando;
    in {
      options.services.vanando = {
        enable = lib.mkEnableOption "Enable vanando image webservice";
        port = lib.mkOption {
          type = lib.types.port;
          default = 8080;
          description = "Port to listen on";
        };
      };

      config = lib.mkIf cfg.enable {
        users = {
          groups.vanando = {};
          users.vanando = {
            isSystemUser = true;
            group = "vanando";
          };
        };
        systemd = {
          sockets.vanando= {
            after = [ "network.target" ];
            wantedBy = [ "sockets.target" ];
            listenStreams = [ "0.0.0.0:${toString cfg.port}" ];
            socketConfig = {
              Accept = false;
              SocketUser = "vanando";
              SocketGroup = "vanando";
              ReusePort = true;
            };
          };

          services.vanando = {
            description = "Vanando image webservice";
            requires = [ "vanando.socket" ];
            serviceConfig = {
              ExecStart = "${vanando}/bin/test";
              Type = "simple";
              PrivateNetwork = true;
              StateDirectory = "vanando";
              User = "vanando";
              Group = "vanando";
              Environment = [
                "PORT=${toString cfg.port}"
              ];
            };
          };

          timers.vanando-scraper = {
            description = "Runs the scraper every hour";
            wantedBy = [ "timers.target" ];

            timerConfig = {
              OnCalendar = "Mon 18:35";
              Persistent = true;
              Unit = "vanando-scraper.service";
            };
          };

          services.vanando-scraper = {
            description = "Vanando image scraper";
            serviceConfig = {
              ExecStart = "${vanando}/bin/scraper";
              Type = "oneshot";
              User = "vanando";
              Group = "vanando";
              StateDirectory = "vanando";
              RemainAfterExit = false;
            };
          };
        };
      };
    };
  };
}

