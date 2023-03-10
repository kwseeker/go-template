# gorm

## 依赖

```shell
go get -u gorm.io/gorm
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



## 参考

+ [GORM 指南](https://gorm.io/zh_CN/docs/index.html)

