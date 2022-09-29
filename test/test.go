package main

import (
	"context"
	"log"

	"github.com/jannchie/gazer-system/api"
	"github.com/jannchie/gazer-system/pkg/client/gs"
)

func main() {

	ctx := context.Background()

	// 定义一个 Worker Group，连接服务器
	wg := gs.NewWorkerGroup([]string{"server:2001"})

	// 定义一个既能解析，又能发起爬取任务的 Worker
	bwu := gs.NewBothWorker(wg.Client, "test", func(tasks chan<- *api.Task) {
		tasks <- &api.Task{
			Url:        "https://pv.sohu.com/cityjson",
			Tag:        "test",
			IntervalMS: 2000, // 表示只爬一次
		}
	}, func(r *api.Raw, c *gs.Client) error {
		log.Printf("%+v\n", r)
		return nil
	})

	// 添加这个 Worker
	wg.AddByWorkUnit(bwu)
	// 跑
	wg.Run(ctx)
}
