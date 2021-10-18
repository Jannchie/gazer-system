package gs

import (
	"context"
	"github.com/jannchie/gazer-system/api"
	"google.golang.org/grpc"
)

type Client struct {
	api.GazerSystemClient
}

func (c *Client) SendOneTask(task *api.Task) error {
	_, err := c.AddTasks(context.Background(), &api.AddTasksReq{Tasks: []*api.Task{task}})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) SendTasks(tasks []*api.Task) error {
	_, err := c.AddTasks(context.Background(), &api.AddTasksReq{Tasks: tasks})
	if err != nil {
		return err
	}
	return nil
}

func NewClient(server string) *Client {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial(server, opts...)
	if err != nil {
		panic(err)
	}
	return &Client{GazerSystemClient: api.NewGazerSystemClient(conn)}
}
