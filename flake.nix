{
  description = "Convert SSH Ed25519 keys to age keys";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs =
    { self, nixpkgs }:
    let
      inherit (nixpkgs) lib;
      systems = [
        "aarch64-linux"
        "x86_64-linux"
        "riscv64-linux"

        "x86_64-darwin"
        "aarch64-darwin"
      ];
      eachSystem = f: lib.genAttrs systems (system: f nixpkgs.legacyPackages.${system});
    in
    {
      packages = eachSystem (pkgs: rec {
        ssh-to-age = pkgs.callPackage ./default.nix { };
        default = ssh-to-age;
      });

      checks = eachSystem (
        pkgs:
        let
          system = pkgs.stdenv.hostPlatform.system;
          packages = lib.mapAttrs' (n: lib.nameValuePair "package-${n}") self.packages.${system};
          devShells = lib.mapAttrs' (n: lib.nameValuePair "devShell-${n}") (self.devShells.${system} or { });
        in
        {
          cross-build = pkgs.callPackage ./cross-build.nix {
            inherit (self.packages.${system}) ssh-to-age;
          };
        }
        // packages
        // devShells
      );
    };
}
