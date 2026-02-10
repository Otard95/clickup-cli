{
  description = "clickup - CLI for interacting with ClickUp";
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  inputs.systems.url = "github:nix-systems/default";
  inputs.flake-utils = {
    url = "github:numtide/flake-utils";
    inputs.systems.follows = "systems";
  };

  outputs =
    { nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go
            cobra-cli
          ];
        };

        packages.default = pkgs.buildGoModule {
          pname = "clickup";
          version = "0.1.0";

          src = ./.;

          vendorHash = null; # will need updating after first `go mod vendor` or nix build

          subPackages = [ "." ];

          meta = with pkgs.lib; {
            description = "CLI for interacting with ClickUp";
            license = licenses.mit;
            mainProgram = "clickup-cli";
          };
        };
      }
    );
}
