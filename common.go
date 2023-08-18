package lib

import (
	"context"
	_ "embed"
	"net"
)

type Resolver interface {
	LookupTXT(context.Context, string) ([]string, error)
}

var resolver Resolver = net.DefaultResolver

func UseResolver(r Resolver) Resolver {
	if r != nil {
		resolver = r
	}

	return resolver
}
