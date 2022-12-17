package lib

import (
	jwt "github.com/dgrijalva/jwt-go"
	mh "github.com/ipsn/go-ipfs/gxlibs/github.com/multiformats/go-multihash"
)

func (claims *Claims) Calculate() {
	hash, _ := mh.Sum([]byte(claims.fqdn+claims.project), mh.SHA1, -1)
	claims.Address = hash.B58String()
}

func (claims *Claims) Valid() error {
	return nil
}

func (claims *Claims) Sign() (Token, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	out, err := token.SignedString(claims.privateKey)
	return Token(out), err
}

func FromToken(token Token, options ...Option) (*Claims, error) {
	claims, err := _new(options)
	if err != nil {
		return nil, err
	}

	_, err = jwt.ParseWithClaims(string(token), claims, func(token *jwt.Token) (interface{}, error) {
		return claims.publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	return claims, nil
}
