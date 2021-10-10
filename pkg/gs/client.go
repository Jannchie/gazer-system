package gs

import (
	"github.com/jannchie/gazer-system/api"
	"google.golang.org/grpc"
)

func NewClient(server string) api.GazerSystemClient {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial(server, opts...)
	if err != nil {
		panic(err)
	}
	return api.NewGazerSystemClient(conn)
}
