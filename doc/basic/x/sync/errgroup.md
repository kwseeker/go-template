# errgroup

errgroup 包为处理公共任务的子任务的 goroutine 组提供了同步、错误传播和上下文取消的功能。

同步并发任务意味着要么等待所有任务完成后再做其他事情，要么在出现问题时取消所有任务。

**实现原理**：



结构体:

```go
type token struct{}

// Group 是执行子任务的go协程的集合
// 空的Group也是有效的，没有限制活跃的goroutines数量，不会在出现错误的时候被取消
type Group struct {
    cancel func()	//ctx取消函数，在WithContext()方法（注意是errgroup的方法）中将ctx.WithCancel返回的cancel函数复制给这里的cancel
	wg sync.WaitGroup
	sem chan token
	errOnce sync.Once
	err     error
}
```

公共方法：

```go
// 创建一个新的Group 以及 ctx驱动的可取消的Context
func WithContext(ctx context.Context) (*Group, context.Context)

```

