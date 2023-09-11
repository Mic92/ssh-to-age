{
  description = "Convert SSH Ed25519 keys to age keys";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-parts.url = "github:/hercules-ci/flake-parts";
  };

  outputs = inputs@{ self, flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } ({ lib, ... }: {
      systems = lib.systems.flakeExposed;
      perSystem = { config, self', inputs', pkgs, system, ... }: {
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
            cross-build = self'.packages.ssh-to-age.overrideAttrs (old: {
              nativeBuildInputs = old.nativeBuildInputs ++ [ pkgs.gox ];
              buildPhase = ''
                runHook preBuild
                HOME=$TMPDIR gox -verbose -osarch '!darwin/386' ./cmd/ssh-to-age/
                runHook postBuild
              '';
              doCheck = false;
              installPhase = ''
                runHook preBuild
                touch $out
                runHook postBuild
              '';
            });
          } // packages // devShells;
      };
    });
}
