package gs

import (
	"log"

	"google.golang.org/grpc/resolver"
)

// Following is an example name resolver. It includes a
// ResolverBuilder(https://godoc.org/google.golang.org/grpc/resolver#Builder)
// and a Resolver(https://godoc.org/google.golang.org/grpc/resolver#Resolver).
//
// A ResolverBuilder is registered for a scheme (in this example, "example" is
// the scheme). When a ClientConn is created for this scheme, the
// ResolverBuilder will be picked to build a Resolver. Note that a new Resolver
// is built for each ClientConn. The Resolver will watch the updates for the
// target, and send updates to the ClientConn.

// GazerSystemResolverBuilder
// ResolverBuilder(https://godoc.org/google.golang.org/grpc/resolver#Builder).
type GazerSystemResolverBuilder struct {
	servers []string
}

func (b *GazerSystemResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, _ resolver.BuildOptions) (resolver.Resolver, error) {
	r := &GazerSystemResolver{
		target: target,
		cc:     cc,
		addressMap: map[string][]string{
			"gazer-system": b.servers,
		},
	}
	r.ResolveNow(resolver.ResolveNowOptions{})
	return r, nil
}
func (*GazerSystemResolverBuilder) Scheme() string { return "gs" }

// GazerSystemResolver
// Resolver(https://godoc.org/google.golang.org/grpc/resolver#Resolver).
type GazerSystemResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addressMap map[string][]string
}

func (r *GazerSystemResolver) ResolveNow(_ resolver.ResolveNowOptions) {
	addrStrings := r.addressMap[r.target.Endpoint]
	addresses := make([]resolver.Address, len(addrStrings))
	for i, s := range addrStrings {
		addresses[i] = resolver.Address{Addr: s}
	}
	err := r.cc.UpdateState(resolver.State{Addresses: addresses})
	if err != nil {
		log.Println(err)
		return
	}
}

func (*GazerSystemResolver) Close() {}

func RegisterServers(servers []string) {
	builder := &GazerSystemResolverBuilder{servers: servers}
	resolver.Register(builder)
}
