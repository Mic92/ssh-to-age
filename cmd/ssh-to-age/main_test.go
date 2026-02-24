package main

import (
	"bufio"
	"fmt"
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
	tempdir, err := os.MkdirTemp(os.TempDir(), "testdir")
	ok(t, err)
	return tempdir
}

func TestPublicKey(t *testing.T) {
	tempdir := TempDir(t)
	defer os.RemoveAll(tempdir) //nolint:errcheck
	out := path.Join(tempdir, "out")

	err := convertKeys([]string{"ssh-to-age", "-i", Asset("id_ed25519.pub"), "-o", out})
	ok(t, err)

	rawPublicKey, err := os.ReadFile(out)
	ok(t, err)
	pubKey := strings.TrimSuffix(string(rawPublicKey), "\n")

	fmt.Printf("public key: %s\n", pubKey)
	_, err = age.ParseX25519Recipient(pubKey)
	ok(t, err)
}

func TestSshKeyScan(t *testing.T) {
	tempdir := TempDir(t)
	defer os.RemoveAll(tempdir) //nolint:errcheck
	out := path.Join(tempdir, "out")

	err := convertKeys([]string{"ssh-to-age", "-i", Asset("keyscan.txt"), "-o", out})
	ok(t, err)

	file, err := os.Open(out)
	ok(t, err)
	defer file.Close() //nolint:errcheck

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
	defer os.RemoveAll(tempdir) //nolint:errcheck
	out := path.Join(tempdir, "out")

	err := convertKeys([]string{"ssh-to-age", "-private-key", "-i", Asset("id_ed25519"), "-o", out})
	ok(t, err)

	rawPrivateKey, err := os.ReadFile(out)
	privateKey := strings.TrimSuffix(string(rawPrivateKey), "\n")
	ok(t, err)

	fmt.Printf("private key: %s\n", privateKey)
	_, err = age.ParseX25519Identity(privateKey)
	ok(t, err)
}

func TestPrivateKeyWithPassphrase(t *testing.T) {
	tempdir := TempDir(t)
	defer os.RemoveAll(tempdir) //nolint:errcheck
	out := path.Join(tempdir, "out")

	passphrase := "test"

	_ = os.Setenv("SSH_TO_AGE_PASSPHRASE", passphrase)
	defer os.Unsetenv("SSH_TO_AGE_PASSPHRASE") //nolint:errcheck

	err := convertKeys([]string{"ssh-to-age", "-private-key", "-i", Asset("id_ed25519_passphrase"), "-o", out})
	ok(t, err)

	rawPrivateKey, err := os.ReadFile(out)
	privateKey := strings.TrimSuffix(string(rawPrivateKey), "\n")
	ok(t, err)

	fmt.Printf("private key: %s\n", privateKey)
	_, err = age.ParseX25519Identity(privateKey)
	ok(t, err)
}

func TestPrivateKeyWithStdinPassphrase(t *testing.T) {
	tempdir := TempDir(t)
	defer os.RemoveAll(tempdir)
	out := path.Join(tempdir, "out")

	stdin, err := os.CreateTemp(tempdir, "stdin")
	ok(t, err)
	defer stdin.Close()

	_, err = stdin.WriteString("test\n")
	ok(t, err)

	_, err = stdin.Seek(0, 0)
	ok(t, err)

	oldStdin := os.Stdin
	os.Stdin = stdin
	defer func() {
		os.Stdin = oldStdin
	}()

	err = convertKeys([]string{"ssh-to-age", "-private-key", "-stdinpass", "-i", Asset("id_ed25519_passphrase"), "-o", out})
	ok(t, err)

	rawPrivateKey, err := os.ReadFile(out)
	privateKey := strings.TrimSuffix(string(rawPrivateKey), "\n")
	ok(t, err)

	fmt.Printf("private key: %s\n", privateKey)
	_, err = age.ParseX25519Identity(privateKey)
	ok(t, err)
}

func TestStdinPassphraseRequiresFileInput(t *testing.T) {
	stdin, err := os.CreateTemp(os.TempDir(), "stdin")
	ok(t, err)
	defer os.Remove(stdin.Name())
	defer stdin.Close()

	oldStdin := os.Stdin
	os.Stdin = stdin
	defer func() {
		os.Stdin = oldStdin
	}()

	err = convertKeys([]string{"ssh-to-age", "-private-key", "-stdinpass"})
	if err == nil {
		t.Fatal("expected error when reading both private key and passphrase from stdin")
	}
	if got, want := err.Error(), "cannot read both private key and passphrase from stdin"; got != want {
		t.Fatalf("unexpected error: got %q, want %q", got, want)
	}
}

func TestVersionFlag(t *testing.T) {
	tempdir := TempDir(t)
	defer os.RemoveAll(tempdir) //nolint:errcheck
	out := path.Join(tempdir, "out")

	err := convertKeys([]string{"ssh-to-age", "-version", "-o", out})
	ok(t, err)

	// Verify that no output file was created when version flag is used
	_, err = os.Stat(out)
	if !os.IsNotExist(err) {
		t.Errorf("output file should not exist when -version flag is used")
	}
}
