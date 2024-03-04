{ pkgs ? import <nixpkgs> {}, vendorHash ? "sha256-yhlbo/ou/o71Wl0ZKb8T4w+0Hhk0BmMQZGmL/ioD6yQ=" }:
pkgs.buildGoModule {
  pname = "ssh-to-age";
  version = "1.1.7";

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
