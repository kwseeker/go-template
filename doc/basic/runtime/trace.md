# go tool trace

**go tool trace 功能**：

+ View trace：查看跟踪
+ Goroutine analysis：Goroutine 分析
+ Network blocking profile：网络阻塞概况
+ Synchronization blocking profile：同步阻塞概况
+ Syscall blocking profile：系统调用阻塞概况
+ Scheduler latency profile：调度延迟概况
+ User defined tasks：用户自定义任务
+ User defined regions：用户自定义区域
+ Minimum mutator utilization：最低 Mutator 利用率

**Demo演示**：

```shell
go run trace_example.go 2> trace.out	#将标准错误信息写入trace.out， 注意 2和>之间不能有空格
go tool trace trace.out					#启动一个web应用展示trace信息
```

> pprof的采样数据也可以输出到文件给trace分析。



## 参考

+ [Go 大杀器之跟踪剖析 trace](https://eddycjy.gitbook.io/golang/di-9-ke-gong-ju/go-tool-trace)