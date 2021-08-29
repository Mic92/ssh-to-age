package agessh

import (
	"crypto/ed25519"
	"crypto/sha512"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"filippo.io/edwards25519"
	"github.com/Mic92/ssh-to-age/bech32"
	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/ssh"
)

type UnsupportedKeyType struct{}

func (UnsupportedKeyType) Error() string {
	return "only ed25519 keys are supported"
}

func ed25519PrivateKeyToCurve25519(pk ed25519.PrivateKey) ([]byte, error) {
	h := sha512.New()
	_, err := h.Write(pk.Seed())
	if err != nil {
		return []byte{}, err
	}
	out := h.Sum(nil)
	return out[:curve25519.ScalarSize], nil
}

func ed25519PublicKeyToCurve25519(pk ed25519.PublicKey) ([]byte, error) {
	// See https://blog.filippo.io/using-ed25519-keys-for-encryption and
	// https://pkg.go.dev/filippo.io/edwards25519#Point.BytesMontgomery.
	p, err := new(edwards25519.Point).SetBytes(pk)
	if err != nil {
		return nil, err
	}
	return p.BytesMontgomery(), nil
}

func SSHPrivateKeyToAge(sshKey []byte) (*string, error) {
	privateKey, err := ssh.ParseRawPrivateKey(sshKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ssh private key: %v", err)
	}

	ed25519Key, ok := privateKey.(*ed25519.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("got %s key type but: %v", reflect.TypeOf(privateKey), UnsupportedKeyType {})
	}
	bytes, err := ed25519PrivateKeyToCurve25519(*ed25519Key)
	if err != nil {
		return nil, err
	}

	s, err := bech32.Encode("AGE-SECRET-KEY-", bytes)
	if err != nil {
		return nil, err
	}
	s = strings.ToUpper(s)
	return &s, nil
}

func SSHPublicKeyToAge(sshKey []byte) (*string, error) {
	var err error
	var pk ssh.PublicKey
	if strings.HasPrefix(string(sshKey), "ssh-") {
		pk, _, _, _, err = ssh.ParseAuthorizedKey(sshKey)
	} else {
		_, _, pk, _, _, err = ssh.ParseKnownHosts(sshKey)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to parse ssh public key: %v", err)
	}
	// We only care about ed25519
	if pk.Type() != ssh.KeyAlgoED25519 {
		return nil, fmt.Errorf("got %s key type but: %v", pk.Type(), UnsupportedKeyType {})
	}
	// Get the bytes
	cpk, ok := pk.(ssh.CryptoPublicKey)
	if !ok {
		return nil, errors.New("BUG! public key does not implement ssh.CryptoPublicKey")
	}
	epk, ok := cpk.CryptoPublicKey().(ed25519.PublicKey)
	if !ok {
		return nil, errors.New("BUG! public key is not of type ed25519.PublicKey")
	}
	// Convert the key to curve ed25519
	mpk, err := ed25519PublicKeyToCurve25519(epk)
	if err != nil {
		return nil, fmt.Errorf("cannot convert ed25519 public key to curve25519: %v", err)
	}
	// Encode the key
	s, err := bech32.Encode("age", mpk)
	if err != nil {
		return nil, fmt.Errorf("cannot encode key as bech32: %v", err)
	}
	return &s, nil
}
