{ pkgs ? import <nixpkgs> {}, vendorHash ? "sha256-MY9txM9SYyQBrm3dIlIwZ9C1Twjp+pDKBHNijGf3mPM=" }:
pkgs.buildGoModule {
  pname = "ssh-to-age";
  version = "1.1.6";

  src = ./.;

  inherit vendorHash;

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
