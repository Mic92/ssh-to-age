{ pkgs ? import <nixpkgs> {}, vendorHash ? "sha256-w/REcFeH58DTQwgotxSBVR4y7aQ9rBDX2U0A4vJno7s=" }:
pkgs.buildGoModule {
  pname = "ssh-to-age";
  version = "1.1.8";

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
