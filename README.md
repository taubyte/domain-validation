# taubyte/domain-validation

[![Release](https://img.shields.io/github/release/taubyte/domain-validation.svg)](https://github.com/taubyte/domain-validation/releases)
[![License](https://img.shields.io/github/license/taubyte/domain-validation)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/taubyte/domain-validation)](https://goreportcard.com/report/taubyte/domain-validation)
[![GoDoc](https://godoc.org/github.com/taubyte/domain-validation?status.svg)](https://pkg.go.dev/github.com/taubyte/domain-validation)
[![Discord](https://img.shields.io/discord/973677117722202152?color=%235865f2&label=discord)](https://tau.link/discord)

Used to validate ownership of DNS domains

## Generate Keys
To use this package you need to generate a symetric pair for keys:

```
openssl ecparam -name prime256v1 -genkey -noout -out private.key
```
```
openssl ec -in tprivate.key -pubout -out public.pem
```

NOTE: the private key need to be kept safe.


## Maintainers
 - Samy Fodil @samyfodil
 - Tafseer Khan @tafseer-khan
 - Aron Jalbuena @arontaubyte
