{
  description = "clickup-cli - CLI for interacting with ClickUp";
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  inputs.systems.url = "github:nix-systems/default";
  inputs.flake-utils = {
    url = "github:numtide/flake-utils";
    inputs.systems.follows = "systems";
  };

  outputs =
    { self, nixpkgs, flake-utils, ... }:
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
          pname = "clickup-cli";
          version = "0.1.0";

          src = self;

          vendorHash = "sha256-7K17JaXFsjf163g5PXCb5ng2gYdotnZ2IDKk8KFjNj0=";

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
