package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"filippo.io/age"
)

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

func Asset(name string) string {
	assets := os.Getenv("TEST_ASSETS")
	if assets == "" {
		assets = "test-assets"
	}
	return path.Join(assets, name)
}

func TempDir(t *testing.T) string {
	tempdir, err := ioutil.TempDir(os.TempDir(), "testdir")
	ok(t, err)
	return tempdir
}

func TestPublicKey(t *testing.T) {
	tempdir := TempDir(t)
	defer os.RemoveAll(tempdir)
	out := path.Join(tempdir, "out")

	err := convertKeys([]string{"ssh-to-age", "-i", Asset("id_ed25519.pub"), "-o", out})
	ok(t, err)

	rawPublicKey, err := ioutil.ReadFile(out)
	ok(t, err)
	pubKey := strings.TrimSuffix(string(rawPublicKey), "\n")

	fmt.Printf("public key: %s\n", pubKey)
	_, err = age.ParseX25519Recipient(pubKey)
	ok(t, err)
}

func TestSshKeyScan(t *testing.T) {
	tempdir := TempDir(t)
	defer os.RemoveAll(tempdir)
	out := path.Join(tempdir, "out")

	err := convertKeys([]string{"ssh-to-age", "-i", Asset("keyscan.txt"), "-o", out})
	ok(t, err)

	file, err := os.Open(out)
	ok(t, err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pubKey := strings.TrimSuffix(scanner.Text(), "\n")
		fmt.Printf("scanned key: %s\n", pubKey)
		_, err = age.ParseX25519Recipient(pubKey)
		ok(t, err)
	}
	ok(t, scanner.Err())
}

func TestPrivateKey(t *testing.T) {
	tempdir := TempDir(t)
	defer os.RemoveAll(tempdir)
	out := path.Join(tempdir, "out")

	err := convertKeys([]string{"ssh-to-age", "-private-key", "-i", Asset("id_ed25519"), "-o", out})
	ok(t, err)

	rawPrivateKey, err := ioutil.ReadFile(out)
	privateKey := strings.TrimSuffix(string(rawPrivateKey), "\n")
	ok(t, err)

	fmt.Printf("private key: %s\n", privateKey)
	_, err = age.ParseX25519Identity(privateKey)
	ok(t, err)
}

func TestPrivateKeyWithPassphrase(t *testing.T) {
	tempdir := TempDir(t)
	defer os.RemoveAll(tempdir)
	out := path.Join(tempdir, "out")

	passphrase := "test"

	os.Setenv("SSH_TO_AGE_PASSPHRASE", passphrase)
	defer os.Unsetenv("SSH_TO_AGE_PASSPHRASE")

	err := convertKeys([]string{"ssh-to-age", "-private-key", "-i", Asset("id_ed25519_passphrase"), "-o", out})
	ok(t, err)

	rawPrivateKey, err := ioutil.ReadFile(out)
	privateKey := strings.TrimSuffix(string(rawPrivateKey), "\n")
	ok(t, err)

	fmt.Printf("private key: %s\n", privateKey)
	_, err = age.ParseX25519Identity(privateKey)
	ok(t, err)
}
