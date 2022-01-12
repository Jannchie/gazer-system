package main

import (
	"context"
	"log"

	"github.com/jannchie/gazer-system/api"
	"github.com/jannchie/gazer-system/pkg/gs"
)

func main() {
	log.SetFlags(log.Lshortfile)
	cli := gs.NewClient(":4747")
	_ = gs.NewClientWithLB(":4747", ":4748", ":4749")
	_, err := cli.ListRaws(context.Background(), &api.ListRawsReq{Tag: "Video's Tags", Limit: 1})
	if err != nil {
		log.Println(err)
	}
	wg := gs.NewWorkerGroup([]string{":4747"}, gs.WithDebug(true), gs.WithSpeedometerServer(""))
	wg.Run(context.Background())
}
