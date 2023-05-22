{
  description = "Flake utils demo";

  inputs = {
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        overlays = [
        ];

        pkgs = import nixpkgs {
          inherit overlays system;
        };
      in
      {
        packages = rec {
          default = circlecicli;
          circlecicli = pkgs.buildGoModule {
            pname = "circlecicli";
            version = "unstable";

            src = ./.;

            CGO_ENABLED = "0";

            vendorHash = "sha256-m5OTb9ugmAczZM/NYH6vnwpDEKVKefOAZY8xAWtezaw=";
          };

        };
        devShells = rec {
          default = empty;

          empty = pkgs.mkShell {
            buildInputs = with pkgs; [
              go
            ];

            CGO_ENABLED = "0";
          };
        };
      }
    );
}
