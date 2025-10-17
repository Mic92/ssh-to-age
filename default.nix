{ pkgs ? import <nixpkgs> {}, vendorHash ? "sha256-aAWyR6f807NXU40Gqfy7567sU89aOIO91xgnQABDs3k=" }:
let
  fs = pkgs.lib.fileset;
in
pkgs.buildGoModule {
  pname = "ssh-to-age";
  version = "1.2.0";

  src = fs.toSource {
    root = ./.;
    fileset = fs.unions [
      ./go.mod
      ./go.sum
      ./cmd
      ./bech32
      ./convert.go
    ];
  };

  inherit vendorHash;

  subPackages = [ "cmd/ssh-to-age" ];

  ldflags = [ "-s" "-w" "-X main.version=${version}" ];

  # golangci-lint is marked as broken on macOS
  nativeBuildInputs = pkgs.lib.optional (!pkgs.stdenv.isDarwin) [ pkgs.golangci-lint ];

  checkPhase = ''
    runHook preCheck
    go test ./...
    runHook postCheck
  '';

  shellHook = ''
    unset GOFLAGS
  '';

  doCheck = true;

  meta = with pkgs.lib; {
    description = "Convert ssh private keys in ed25519 format to age keys";
    homepage = "https://github.com/Mic92/ssh-to-age";
    license = licenses.mit;
    maintainers = with maintainers; [ mic92 ];
  };
}
