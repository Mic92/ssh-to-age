{
  description = "Convert SSH Ed25519 keys to age keys";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-parts.url = "github:/hercules-ci/flake-parts";
    flake-parts.inputs.nixpkgs-lib.follows = "nixpkgs";
  };

  outputs =
    inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } (
      { lib, ... }:
      {
        systems = [
          "aarch64-linux"
          "x86_64-linux"
          "riscv64-linux"

          "x86_64-darwin"
          "aarch64-darwin"
        ];
        perSystem =
          {
            config,
            self',
            pkgs,
            ...
          }:
          {
            packages = {
              ssh-to-age = (pkgs.callPackage ./default.nix { });
              default = config.packages.ssh-to-age;
            };
            checks =
              let
                packages = lib.mapAttrs' (n: lib.nameValuePair "package-${n}") self'.packages;
                devShells = lib.mapAttrs' (n: lib.nameValuePair "devShell-${n}") self'.devShells;
              in
              {
                cross-build = pkgs.callPackage ./cross-build.nix { 
                  ssh-to-age = self'.packages.ssh-to-age;
                };
              }
              // packages
              // devShells;
          };
      }
    );
}
