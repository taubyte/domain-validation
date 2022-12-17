package lib

import "crypto/ecdsa"

type Token string

type Repository struct {
	Id string `json:"id"`
}

type Claims struct {
	fqdn       string
	project    string
	Address    string `json:"address"` // address is hash(fqdn+projectID)
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}
