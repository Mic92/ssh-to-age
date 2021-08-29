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

- Exports the public key:

```console
$ ssh-to-age -i $HOME/.ssh/id_ed25519.pub -o pub-key.txt
$ cat pub-key.txt
age17044m9wgakla6pzftf4srtl3h5mcsr85jysgt5fg23zpnta8sfdqhzn452
```

ssh-to-age also supports multiple public keys at once seperated by newlines and ignores unless ssh keys that are not in the ed25519 format. This makes it suiteable in combination with `ssh-keyscan`:

```
ssh-keyscan eve.thalheim.io eva.thalheim.io | ssh-to-age
```

## Install with nix

```console
$ nix-shell -p 'import (fetchTarball "https://github.com/Mic92/ssh-to-age/archive/main.tar.gz") {}'
```

## Install with go

```console
$ go get github.com/Mic92/ssh-to-age/cmd
```
