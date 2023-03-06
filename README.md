# ssh-to-age
Convert SSH Ed25519 keys to [age](https://github.com/FiloSottile/age) keys.
This is useful for usage in [sops-nix](https://github.com/Mic92/sops-nix) and
[sops](https://github.com/mozilla/sops)

## Usage

- Exports the private key:

```console
$ ssh-to-age -private-key -i $HOME/.ssh/id_ed25519 -o key.txt
$ cat key.txt
 AGE-SECRET-KEY-1K3VN4N03PTHJWSJSCCMQCN33RY5FSKQPJ4KRRTG3JMQUYE0TUSEQEDH6V8
```

If you private key is encrypted, you can export the password in `SSH_TO_AGE_PASSPHRASE`

``` console
$ read -s SSH_TO_AGE_PASSPHRASE; export SSH_TO_AGE_PASSPHRASE
$ ssh-to-age -private-key -i $HOME/.ssh/id_ed25519 -o key.txt
```

- Exports the public key:

```console
$ ssh-to-age -i $HOME/.ssh/id_ed25519.pub -o pub-key.txt
$ cat pub-key.txt
age17044m9wgakla6pzftf4srtl3h5mcsr85jysgt5fg23zpnta8sfdqhzn452
```

ssh-to-age also supports multiple public keys at once seperated by newlines and ignores unless ssh keys that are not in the ed25519 format. This makes it suiteable in combination with `ssh-keyscan`:

```console
$ ssh-keyscan eve.thalheim.io eva.thalheim.io | ssh-to-age
# eve.thalheim.io:22 SSH-2.0-OpenSSH_8.6
...
age1hjm3aujg9e79f5yth8a2cejzdjg5n9vnu96l05p70uvfpeltnpms7yy3pp
age1v8zjc47jmlqwefyu66s0d4ke98qr4vnuj3cpvs4z9npfdw833dxqwjrhzv
```

## Install with nix

```console
$ nix-shell -p 'import (fetchTarball "https://github.com/Mic92/ssh-to-age/archive/main.tar.gz") {}'
```

## Install with go

```console
$ go install github.com/Mic92/ssh-to-age/cmd/ssh-to-age@latest
```
