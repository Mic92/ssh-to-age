module github.com/Mic92/ssh-to-age

go 1.23.0 // tagx:compat 1.16

toolchain go1.24.3

require (
	filippo.io/age v1.2.1
	filippo.io/edwards25519 v1.1.0
	golang.org/x/crypto v0.37.0
)

require golang.org/x/sys v0.32.0 // indirect
