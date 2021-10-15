package main

import (
	"context"
	"github.com/jannchie/gazer-system/api"
	"github.com/jannchie/gazer-system/pkg/gs"
	"log"
)

func main() {
	ctx := context.Background()
	wg := gs.NewWorkerGroup("localhost:2000")
	bwu := gs.NewBothWorker(wg.Client, "test", func(tasks chan<- *api.Task) {
		tasks <- &api.Task{
			Url:        "https://pv.sohu.com/cityjson",
			Tag:        "test",
			IntervalMS: 0,
		}
	}, func(raw *api.Raw) error {
		log.Printf("%+v\n", raw)
		return nil
	})

	wg.AddByWorkUnit(bwu)

	wg.AddByWorkUnit(gs.NewParserWorker(wg.Client, "?", func(raw *api.Raw) error {
		log.Println(raw)
		return nil
	}))

	wg.Run(ctx)
}
