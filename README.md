# go-zero实战：让微服务Go起来
这是一个 `go-zero` 入门学习教程的 `DTM` 分布式事务示例代码，教程地址：[go-zero实战：让微服务Go起来](https://juejin.cn/post/7036011047391592485)。

#### 子事务屏障相关的表
```sql
create database if not exists dtm_barrier
/*!40100 DEFAULT CHARACTER SET utf8mb4 */
;
drop table if exists dtm_barrier.barrier;
create table if not exists dtm_barrier.barrier(
  id bigint(22) PRIMARY KEY AUTO_INCREMENT,
  trans_type varchar(45) default '',
  gid varchar(128) default '',
  branch_id varchar(128) default '',
  op varchar(45) default '',
  barrier_id varchar(45) default '',
  reason varchar(45) default '' comment 'the branch type who insert this record',
  create_time datetime DEFAULT now(),
  update_time datetime DEFAULT now(),
  key(create_time),
  key(update_time),
  UNIQUE key(gid, branch_id, op, barrier_id)
);
```

## 感谢

- [go-zero](https://github.com/zeromicro/go-zero)
- [DTM](https://github.com/dtm-labs/dtm)