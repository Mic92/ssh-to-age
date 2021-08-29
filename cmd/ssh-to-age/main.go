package main

import (
	"errors"
	"flag"
	"fmt"
	sshage "github/Mic92/ssh-to-age"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type options struct {
	out, in    string
	privateKey bool
}

func parseFlags(args []string) options {
	var opts options
	f := flag.NewFlagSet(args[0], flag.ExitOnError)
	f.BoolVar(&opts.privateKey, "private-key", false, "convert private key instead of public key")
	f.StringVar(&opts.in, "i", "-", "Input path. Reads by default from standard output")
	f.StringVar(&opts.out, "o", "-", "Output path. Prints by default to standard output")
	if err := f.Parse(args[1:]); err != nil {
		// should never happen since flag.ExitOnError
		panic(err)
	}

	return opts
}

//func main() {
//	if len(os.Args) != 1 {
//		println("Usage: " + os.Args[0])
//		println("Pipe a SSH public key or the output of ssh-keyscan into it")
//		os.Exit(1)
//	}
//
//	scanner := bufio.NewScanner(os.Stdin)
//	for scanner.Scan() {
//		text := scanner.Text() + "\n"
//		var err error
//		fmt.Println(s)
//	}
//	if err := scanner.Err(); err != nil {
//		panic(err)
//	}
//}

func convertKeys(args []string) error {
	opts := parseFlags(args)

	var sshKey []byte
	var err error
	if opts.in == "-" {
		sshKey, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("error reading stdin: %w", err)
		}
	} else {
		sshKey, err = ioutil.ReadFile(opts.in)
		if err != nil {
			return fmt.Errorf("error reading %s: %w", opts.in, err)
		}
	}

	writer := io.WriteCloser(os.Stdout)
	if opts.out != "-" {
		writer, err := os.Create(opts.out)
		if err != nil {
			return fmt.Errorf("failed to create %s: %w", opts.out, err)
		}
		defer writer.Close()
	}

	if opts.privateKey {
		key, err := sshage.SSHPrivateKeyToAge(sshKey)
		if err != nil {
			return fmt.Errorf("failed to convert '%s': %v", sshKey, err)
		}
		if _, err := writer.Write(key); err != nil {
			return fmt.Errorf("failed to write key: %v", err)
		}
	} else {
		keys := strings.Split(string(sshKey), "\n")
		for _, k := range keys {
			key, err := sshage.SSHPrivateKeyToAge([]byte(k))
			if err != nil {
				if errors.Is(err, sshage.UnsupportedKeyType) {
					// silently ignore unsupported key types to make it ssh-keyscan friendly
					continue
				}
				return fmt.Errorf("failed to convert '%s': %v", k, err)
			}
			if _, err := writer.Write(key); err != nil {
				return fmt.Errorf("failed to write key: %v", err)
			}
		}
	}
	return nil
}

func main() {
	if err := convertKeys(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err)
		os.Exit(1)
	}
}
