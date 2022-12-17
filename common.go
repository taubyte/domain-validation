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
	resolver = r
	return resolver
}
