package main

import (
	"context"
	"github.com/jannchie/gazer-system/api"
	"github.com/jannchie/gazer-system/pkg/gs"
	"log"
)

func main() {
	cli := gs.NewClient(":2000")
	_ = gs.NewClientWithLB(":2000", ":2001", ":2002")
	res, err := cli.ListRaws(context.Background(), &api.ListRawsReq{Tag: "Video's Tags", Limit: 1})
	log.Println(res)
	if err != nil {
		log.Println(err)
	}
}
