# go-zero实战：让微服务Go起来
这是一个 `go-zero` 入门学习教程的示例代码，教程地址：[o-zero实战：让微服务Go起来](https://juejin.cn/post/7036011047391592485)。

### docker 本地集成环境安装
```bash
$ cd gonivinck
$ docker-compose up -d
```

### 数据库创建
使用 `mysql` 管理工具，访问 `127.0.0.1:3306`，用户名：`root`，密码：`123456`，创建数据库 `mall`，创建数据表 `user`，`product`，`order`，`pay`，`sql` 语句参见教程或者目录`code/mall`下。

### 项目启动
进入`golang`容器，根据教程启动所有服务。

