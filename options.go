package lib

import (
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/ipfs/go-cid"
)

type Option func(*Claims) error

func FQDN(fqdn string) Option {
	return func(c *Claims) error {
		c.fqdn = fqdn
		return nil
	}
}

func Project(project cid.Cid) Option {
	return func(c *Claims) error {
		c.project = project.String()
		return nil
	}
}

func PrivateKey(data []byte) Option {
	return func(c *Claims) (err error) {
		c.privateKey, err = jwt.ParseECPrivateKeyFromPEM(data)
		if err != nil {
			err = fmt.Errorf("Parsing JWT private key failed with: %w", err)
		}
		return
	}
}

func PublicKey(data []byte) Option {
	return func(c *Claims) (err error) {
		c.publicKey, err = jwt.ParseECPublicKeyFromPEM(data)
		if err != nil {
			err = fmt.Errorf("Parsing JWT public key failed with: %w", err)
		}
		return
	}
}
