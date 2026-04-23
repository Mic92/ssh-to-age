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

If your private key is encrypted, pipe the passphrase via `-stdinpass`:

```console
$ systemd-ask-password | ssh-to-age -private-key -stdinpass -i $HOME/.ssh/id_ed25519 -o key.txt
$ security find-generic-password -w -s 'SSH Key Passphrase' | ssh-to-age -private-key -stdinpass -i $HOME/.ssh/id_ed25519 -o key.txt
```

Alternatively, pass it via `SSH_TO_AGE_PASSPHRASE` environment variable.

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

## Security considerations

`ssh-to-age` deterministically derives an X25519 (age) key from an Ed25519 SSH
key using the standard Ed25519-to-X25519 birational map (see
[Filippo Valsorda's write-up](https://blog.filippo.io/using-ed25519-keys-for-encryption/)):

- The age private scalar is `sha512(seed)[:32]` — the same scalar Ed25519
  itself derives from the seed per RFC 8032. No new key material is generated.
- The age public key is the Montgomery form of the Ed25519 public point.
- The mapping is **deterministic** in one direction only. You cannot recover
  the original SSH private key from an age key: SHA-512 is one-way, and the
  OpenSSH key file (seed, nonce-generation half of the hash, optional
  passphrase-based encryption) cannot be reconstructed from the X25519 scalar.
- The age key does **not** grant SSH authentication — Ed25519 signing also
  requires the raw seed, which is not recoverable from the age key.

### Storing a derived age key on disk

A common pattern is using a passphrase-protected SSH key together with an
unencrypted age key for `sops` / `sops-nix`. Be aware that:

- An attacker with read access to the age key file can decrypt every secret
  encrypted to the derived age public key. This bypasses the passphrase
  protection of the SSH key *for decryption purposes*.
- SSH authentication itself remains protected by the passphrase.

If you want the age key to inherit the at-rest protection of your SSH key,
do **not** derive it from that SSH key. Generate a dedicated age identity and
protect it separately — for example with
[`age-plugin-yubikey`](https://github.com/str4d/age-plugin-yubikey),
[`age-plugin-tpm`](https://github.com/Foxboron/age-plugin-tpm), or by keeping
it on an encrypted volume that is unlocked on demand.

## Install with nix

```console
$ nix-shell -p 'import (fetchTarball "https://github.com/Mic92/ssh-to-age/archive/main.tar.gz") {}'
```

## Install with go

```console
$ go install github.com/Mic92/ssh-to-age/cmd/ssh-to-age@latest
```
