package main

import (
	"context"
	"github.com/jannchie/gazer-system/api"
	"google.golang.org/grpc"
	"log"
)

func main() {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial("localhost:2000", opts...)
	if err != nil {
		panic(err)
	}
	client := api.NewGazerSystemClient(conn)
	ctx := context.Background()

	raws, err := client.ListRaws(ctx, &api.ListRawsReq{
		Tag:   "test",
		Limit: 100,
	})
	if err != nil {
		panic(err)
	}
	var idList []uint64
	for i := range raws.GetRaws() {
		idList = append(idList, raws.GetRaws()[i].Id)
	}
	_, _ = client.ConsumeRaws(ctx, &api.ConsumeRawsReq{IdList: idList})
	res, err := client.AddTasks(ctx, &api.AddTasksReq{
		Tasks: []*api.Task{
			{
				Url:        "https://www.baidu.com",
				Tag:        "test",
				IntervalMS: 0,
			},
		},
	})
	if err != nil {
		panic(err)
	}
	log.Println(res)

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}(conn)
}
