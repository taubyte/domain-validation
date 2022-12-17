package lib

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/ipfs/go-cid"

	mh "github.com/ipsn/go-ipfs/gxlibs/github.com/multiformats/go-multihash"
)

func (t Token) Bytes() ([]byte, error) {
	return hex.DecodeString(string(t))
}

func TXTEntry(project *cid.Cid, fqdn string) string {
	return project.String()[:8] + "." + fqdn
}

func Address(project *cid.Cid, fqdn string) string {
	hash, _ := mh.Sum([]byte(fqdn+project.String()), mh.SHA1, -1)
	return hash.B58String()
}

func FromDNS(ctx context.Context, project *cid.Cid, fqdn string, options ...Option) error {
	txtRecords, err := resolver.LookupTXT(ctx, TXTEntry(project, fqdn))
	if err != nil {
		return fmt.Errorf("TXT lookup failed with: %w", err)
	}

	for _, token := range txtRecords {
		claims, err := FromToken(Token(token), options...)
		if err == nil {
			address := Address(project, fqdn)
			if claims.Address == address {
				return nil
			}
		}
	}

	return errors.New("failed to verify token")
}
