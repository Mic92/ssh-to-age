{
  description = "Convert SSH Ed25519 keys to age keys";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-parts.url = "github:/hercules-ci/flake-parts";
  };

  outputs = inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } ({ lib, ... }: {
      systems = lib.systems.flakeExposed;
      perSystem = { config, self', inputs', pkgs, system, ... }:  {
        packages = {
          ssh-to-age = (pkgs.callPackage ./default.nix {});
          default = config.packages.ssh-to-age;
        };
      };
    });
}
