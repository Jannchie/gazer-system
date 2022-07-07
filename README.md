# Gazer System - 凝视系统

这是一个用于长期追踪数据变化的工具。

我们常需要定时针对某些链接进行爬取，以获取历史数据，这个系统为此而生。这是一组无状态的服务，意味着它可以简单地水平扩展。它提供了如下几个功能：

1. 自动代理：使用洋葱路由进行代理，无需自行维护 IP 池，也可以绕过大多数基于 IP 地址的反爬虫措施。
2. 计划爬取：可以定义爬取计划，以一定间隔进行接口的访问。
3. 并发爬取：使用 Golang 开发，支持高速并发爬取，IO 才是唯一瓶颈。
4. 数据存储：支持 Sqlite 数据存储，无需配置数据库，开箱即用。

## 架构

整个系统采用 CS 架构，服务器能够并发消费爬取任务，并提供原始数据。客户端可以批量拉取原始数据进行进一步解析。
整个过程提供类似消息队列的确认机制和重试机制，防止爬取任务重复或者丢失。

## 服务端部署

建议使用 Docker Compose 部署。否则需要自行配置代理服务器并进行构建（目前没有提供构建步骤）。

使用 Docker Compose 部署，需要提供的 yaml 示例文件如下：

```yaml
version: "3"
volumes:
  db:
services:
  server:
    image: jannchie/gazer-system-server
    ports:
      - "2000:2000"
    volumes:
      - db:/data
    environment:
      TOR: "proxy:9050"
      TOR_CTL: "proxy:9051"
      PORT: 2000
      DSN: "/data/gazer-system.db"
      TOR_PASSWORD: ${TOR_PASSWORD}
  proxy:
    image: dperson/torproxy
    expose:
      - 9050
      - 9051
    environment:
      PASSWORD: ${TOR_PASSWORD}
```

使用如下命令即可启动服务器：

```bash
docker compose pull
docker compose up
```

## 客户端

目前只开发了 Golang 的客户端。

示例程序如下：

``` golang
func main() {
	ctx := context.Background()

	// 定义一个 Worker Group，连接服务器
	wg := gs.NewWorkerGroup([]string{"localhost:2000"})

	// 定义一个既能解析，又能发起爬取任务的 Worker
	bwu := gs.NewBothWorker(wg.Client, "test", func(tasks chan<- *api.Task) {
		tasks <- &api.Task{
			Url:        "https://pv.sohu.com/cityjson",
			Tag:        "test",
			IntervalMS: 0, // 表示只爬一次
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
```
