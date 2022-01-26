# go-zero实战：让微服务Go起来
这是一个 `go-zero` 入门学习教程的示例代码，教程地址：[go-zero实战：让微服务Go起来](https://juejin.cn/post/7036011047391592485)。

`DTM` 分布式事务示例代码请切换至 [dtm](https://github.com/nivin-studio/go-zero-mall/tree/dtm) 分支。

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
下载本项目，将项目代码放置 `gonivinck` 配置 `CODE_PATH_HOST` 指定的本机目录，进入 `golang` 容器，运行项目代码。

#### 3.1 进入 `golang` 容器
~~~bash
$ docker exec -it gonivinck_golang_1 bash
~~~

#### 3.2 使用 `nivin` 命令工具

- nivin install
安装项目依赖命令。

~~~bash
$ ./nivin install
~~~

- nivin start [rpc|api] [service_name]
服务启动命令，创建服务会话，并启动对应的服务。
    
~~~bash
$ ./nivin start rpc user
~~~

~~~bash
$ ./nivin start api user
~~~

- nivin stop [rpc|api] [service_name]
服务暂停命令，删除对应的服务会话。
    
~~~bash
$ ./nivin stop rpc user
~~~

~~~bash
$ ./nivin stop api user
~~~

- nivin info [rpc|api] [service_name]
服务查看命令，可以进入服务对应的会话终端，查看运行日志。

~~~bash
$ ./nivin info rpc user
~~~

~~~bash
$ ./nivin info api user
~~~

> 提示：使用 ctrl+a+d 组合快捷键，可以无损退出此会话，不会中止会话中运行的服务。

- nivin ls
服务会话列表，查看启动的服务会话列表。
    
~~~bash
$ ./nivin ls
~~~


## 感谢

- [go-zero](https://github.com/zeromicro/go-zero)
- [DTM](https://github.com/dtm-labs/dtm)