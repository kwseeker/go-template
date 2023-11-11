# gorm

**特性**：

- 全功能 ORM
- 关联 (Has One，Has Many，Belongs To，Many To Many，多态，单表继承)
- Create，Save，Update，Delete，Find 中钩子方法
- 支持 `Preload`、`Joins` 的预加载
- 事务，嵌套事务，Save Point，Rollback To Saved Point
- Context、预编译模式、DryRun 模式
- 批量插入，FindInBatches，Find/Create with Map，使用 SQL 表达式、Context Valuer 进行 CRUD
- SQL 构建器，Upsert，数据库锁，Optimizer/Index/Comment Hint，命名参数，子查询
- 复合主键，索引，约束
- Auto Migration
- 自定义 Logger
- 灵活的可扩展插件 API：Database Resolver（多数据库，读写分离）、Prometheus…
- 每个特性都经过了测试的重重考验
- 开发者友好



## 使用

### 依赖

```shell
go get -u gorm.io/gorm
# 假如连接MySQL
go get -u gorm.io/driver/mysql
```



## 功能分析

### [迁移](https://gorm.io/zh_CN/docs/migration.html)

最开始看到这个还以为是支持**线上不停服迁移**（线上大表改表结构其实也属于数据迁移问题），然后看了下`AutoMigrate()`方法的源码实现后大失所望，GORM的迁移功能只是用于自动化改表用的。官方原话：

"AutoMigrate 用于自动迁移您的 schema，保持您的 schema 是最新的。"

其实它并不关心改表结构过程中锁表了怎么办。所以千万不能用它来改线上大表的表结构。

`AutoMigrate()`实现原理简单说就是：

1. 克隆一个DB连接；
2. 查询数据库的元数据库看下表对象（如 struct User{}）对应的数据表是否存在，不存在就直接按表对象数据结构创建新表；
3. 存在的话，解析表对象的字段和字段标签和数据表字段元数据对比，查找变更的字段，构造对应的sql进行更新。

比如 `migrator_test.go`测试日志

```txt
=== RUN   TestAutoMigrate
2023-03-10T14:33:27+08:00 error layer=debugger gnu_debuglink CRC check failed for /lib/x86_64-linux-gnu/libc-2.27.so (want 42cc048 got b1c74187)

2023/03/10 14:34:16 /home/lee/go/pkg/mod/gorm.io/driver/mysql@v1.3.6/migrator.go:250
[0.341ms] [rows:-] SELECT DATABASE()

2023/03/10 14:34:16 /home/lee/go/pkg/mod/gorm.io/driver/mysql@v1.3.6/migrator.go:253
[3.203ms] [rows:1] SELECT SCHEMA_NAME from Information_schema.SCHEMATA where SCHEMA_NAME LIKE 'test%' ORDER BY SCHEMA_NAME='test' DESC,SCHEMA_NAME limit 1

2023/03/10 14:34:16 /home/lee/go/src/kwseeker.top/kwseeker/go-template/orm/gorm/migrator/migrator_test.go:36
[1.145ms] [rows:-] SELECT count(*) FROM information_schema.tables WHERE table_schema = 'test' AND table_name = 'users' AND table_type = 'BASE TABLE'

2023/03/10 14:34:28 /home/lee/go/pkg/mod/gorm.io/driver/mysql@v1.3.6/migrator.go:250
[0.370ms] [rows:-] SELECT DATABASE()

2023/03/10 14:34:28 /home/lee/go/pkg/mod/gorm.io/driver/mysql@v1.3.6/migrator.go:253
[1.581ms] [rows:1] SELECT SCHEMA_NAME from Information_schema.SCHEMATA where SCHEMA_NAME LIKE 'test%' ORDER BY SCHEMA_NAME='test' DESC,SCHEMA_NAME limit 1

2023/03/10 14:34:28 /home/lee/go/pkg/mod/gorm.io/driver/mysql@v1.3.6/migrator.go:163
[0.678ms] [rows:-] SELECT * FROM `users` LIMIT 1

2023/03/10 14:34:28 /home/lee/go/pkg/mod/gorm.io/driver/mysql@v1.3.6/migrator.go:181
[2.129ms] [rows:-] SELECT column_name, column_default, is_nullable = 'YES', data_type, character_maximum_length, column_type, column_key, extra, column_comment, numeric_precision, numeric_scale , datetime_precision FROM information_schema.columns WHERE table_schema = 'test' AND table_name = 'users' ORDER BY ORDINAL_POSITION

2023/03/10 14:40:54 /home/lee/go/src/kwseeker.top/kwseeker/go-template/orm/gorm/migrator/migrator_test.go:36
[44.020ms] [rows:0] ALTER TABLE `users` ADD `gender` tinyint DEFAULT 0
2023/03/10 14:42:38 migrate done
```
### 连接

+ 数据源名称DSN
+ 连接配置
  + 默认事务开关
  + 命令策略
  + Logger
  + NowFunc创建事件生成函数
  + DryRun开关
  + 缓存预编译语句开关
  + 嵌套事务开关
  + 启用全局 update/delete ???
  + DisableAutomaticPing 自动 ping 开关
  + 迁移时自动创建外键约束开关
+ 自定义驱动
+ Session ???
+ 连接池
+ 多数据库支持DBResolver
  + 支持多个 sources、replicas
  + 读写分离
  + 根据工作表、struct 自动切换连接
  + 手动切换连接
  + Sources/Replicas 负载均衡
  + 适用于原生 SQL
  + 事务

### CRUD

+ 对象与表映射模型
  + 字段默认值
  + 自定义数据类型

+ 增
  + 使用*对象指针*插入一条记录
    + 指定字段插入
    + 忽略字段插入
  + 使用*对象切片指针*批量插入
    + 一批
    + 分批
  + 插入Hook
  + 使用*Map*插入记录
  + 使用*SQL 表达式＋Context Valuer*插入记录
  + 使用*关联*插入记录
  + Upsert（不存在则插入存在则更新）及冲突
+ 查
  + 查询Hook
+ 改
+ 删

+ 使用原生SQL
+ Scopes复用通用逻辑
+ Set, Get, InstanceSet, InstanceGet 传值
+ 对数据进行序列化和反序列化
+ 建表或迁移建表
  + 创建索引
  + 创建约束
  + 复合主键

### 事务

### 性能

+ 缓存预编译语句

+ Index Hints
+ 读写分离
+ Sharing分表

### 安全

+ 防止注入

### 插件

### 辅助功能

#### 自动迁移

#### Logger

#### Prometheus




## 参考

+ [GORM 指南](https://gorm.io/zh_CN/docs/index.html)

