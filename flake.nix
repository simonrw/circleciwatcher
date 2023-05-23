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
          default = circleciwatcher;
          circleciwatcher = pkgs.buildGoModule {
            pname = "circleciwatcher";
            version = "unstable";

            src = ./.;

            CGO_ENABLED = "0";

            vendorHash = null;
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
