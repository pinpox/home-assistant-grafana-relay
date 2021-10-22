{
  description = "Relay grafana alerts to home-assistant notifications";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    flake-compat = {
      url = "github:edolstra/flake-compat";
      flake = false;
    };
  };

  outputs = { self, nixpkgs, flake-utils, ... }:

    flake-utils.lib.eachDefaultSystem (system:
      let pkgs = nixpkgs.legacyPackages.${system};
      in rec {
        packages = flake-utils.lib.flattenTree rec {

          ha-relay = pkgs.buildGoModule rec {
            pname = "ha-relay";
            version = "1.0.0";
            src = self;
            vendorSha256 = null;

            meta = with pkgs.lib; {
              maintainers = with maintainers; [ pinpox ];
              license = licenses.gpl3;
              description = "Grafana Alerts to home-assistant notifications";
              homepage = "https://github.com/pinpox/home-assistant-grafana-relay";
            };
          };
        };

        apps = {
          ha-relay = flake-utils.lib.mkApp {
            drv = packages.ha-relay;
            exePath = "/bin/home-assistant-grafana-relay";
          };
        };
        defaultPackage = packages.ha-relay;

        defaultApp = apps.ha-relay;
      });
}
