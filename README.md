# Gazer System - 凝视系统

这是一个用于长期追踪数据变化的工具。

我们常需要定时针对某些链接进行爬取，以获取历史数据。这个系统为此而生，它提供了如下几个功能：

1. 自动代理：使用洋葱路由进行代理，无需自行维护 IP 池，也可以绕过大多数基于 IP 地址的反爬虫措施。
2. 并发爬取：使用 Golang 开发，支持高速并发爬取，IO 才是唯一瓶颈。
3. 数据存储：支持 Sqlite 数据存储，无需配置数据库，开箱即用。


