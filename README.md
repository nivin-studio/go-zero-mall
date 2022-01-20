# go-zero实战：让微服务Go起来
这是一个 `go-zero` 入门学习教程的示例代码，教程地址：[go-zero实战：让微服务Go起来](https://juejin.cn/post/7036011047391592485)。

## 使用

### 1. `docker` 本地开发环境安装
下载 [gonivinck](https://github.com/nivin-studio/gonivinck) 本地开发环境.

### 2. 数据库创建
地址：`127.0.0.1:3306`

用户：`root`

密码：`123456`

创建数据库 `mall`

创建数据表 `user`、`product`、`order`、`pay`

`SQL`语句在 `service/[user,product,order,pay]/model` 目录下。

> 提示：如果你修改 gonivinck 相关 mysql 配置，请使用你修改的端口号，账号，密码连接访问数据库。

### 3. 项目启动
下载本项目，将项目代码放置 `gonivinck` 配置 `CODE_PATH_HOST` 指定的本机目录，进入 golang 容器，运动项目代码。
- 进入 `golang` 容器
~~~bash
$ docker exec -it gonivinck_golang_1 bash
~~~
- 容器中启动 `user rpc` 服务
~~~bash
$ cd mall/service/user/rpc
$ go run user.go -f etc/user.yaml
Starting rpc server at 127.0.0.1:9000...
~~~
