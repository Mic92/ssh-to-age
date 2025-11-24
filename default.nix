{ pkgs ? import <nixpkgs> {}, vendorHash ? "sha256-N9NU1AjPkC/cHuRoXwcwRAi8SbzYo1khTX3jC83I6jY=" }:
let
  fs = pkgs.lib.fileset;
  version = "1.2.0";
in
pkgs.buildGoModule {
  pname = "ssh-to-age";
  inherit version;

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
