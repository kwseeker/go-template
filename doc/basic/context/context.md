# context

Go 1.7 标准库引入 context，中文译作“上下文”，准确说它是 goroutine 的上下文，包含 goroutine 的运行状态、环境、现场等信息。

context 主要用来在 goroutine 之间传递上下文信息，包括：取消信号、超时时间、截止时间、k-v 等。

在Go 里，我们不能直接杀死协程，协程的关闭一般会用 `channel+select` 方式来控制。但是在某些场景下，例如处理一个请求衍生了很多协程，这些协程之间是相互关联的：需要共享一些全局变量、有共同的 deadline 等，而且可以同时被关闭。再用 `channel+select` 就会比较麻烦，这时就可以通过 context 来实现。



## context 原理

**Context** 是一个接口，定义了4个方法：

```go
// 返回一个 channel，可以表示 context 被取消的信号：当这个 channel 被关闭时，说明 context 被取消了
// 这是一个只读的channel。读一个关闭的 channel 会读出相应类型的零值。
// 并且源码里没有地方会向这个 channel 里面塞入值。换句话说，这是一个 receive-only 的 channel。
// 因此在子协程里读这个 channel，除非被关闭，否则读不出来任何东西。
Done() <-chan struct{}
// 在 channel Done 关闭后，返回 context 取消原因
Err() error
// 返回 context 是否会被取消以及自动取消时间（即 deadline）
Deadline() (deadline time.Time, ok bool)
// 获取 key 对应的 value
Value(key interface{}) interface{}
```

**canceler** 接口：

```go
// 定义了一种可以取消的Context类型， cancelCtx 和 timerCtx 实现了此接口
type canceler interface {
	cancel(removeFromParent bool, err error)
	Done() <-chan struct{}
}
```

官方提供了4种Context实现：

+ **emptyCtx** 空的Context实现

  ```go
  type emptyCtx int	//甚至都不是struct
  
  var (
  	background = new(emptyCtx)
  	todo       = new(emptyCtx)
  )
  ```

+ **cancelCtx** 用于cancel协程

  内部还是通过`channel+select`实现的：cancel() 关闭 Done()返回的 channel, go 协程检测到channel 关闭后，退出协程。

  通过遍历多叉树父Context cancel 会逐层调用 cancel() 关闭所有子Context。

  ```go
  type cancelCtx struct {
  	Context					//指向父Context的指针
  	mu       sync.Mutex            // protects following fields
  	done     atomic.Value          // of chan struct{}, created lazily, closed by first cancel call
  	children map[canceler]struct{} // 指向子Context的指针Set集合, 配合Context，正好构成多叉树结构
  	err      error                 // set to non-nil by the first cancel call
  }
  ```

  cancelCtx的数据结构就是一个多叉树（借一张图）：

  ![](https://pic4.zhimg.com/v2-f7ea0b0baec68b718a514419636e875b_r.jpg)

  + **timerCtx** 带超时自动cancel的cancelCtx

    如果父Context deadline 在当前timerCtx deadline之前，则使用父Context的deadline作为cancel的截止时间。

    ```go
    type timerCtx struct {
    	cancelCtx			//指向父Context的指针
    	timer *time.Timer 	// Under cancelCtx.mu.
    	deadline time.Time
    }
    
    //cancel传递，parent父Context, child当前Context
    //内嵌结构体（指向父Context的指针） + chi 
    func propagateCancel(parent Context, child canceler) {
    	done := parent.Done()
    	if done == nil {
    		return // parent is never canceled
    	}
    
    	select {
    	case <-done:	//父Context已经canceled
    		child.cancel(false, parent.Err())	//cancel child Context
    		return
    	default:
    	}
    
    	if p, ok := parentCancelCtx(parent); ok {
    		p.mu.Lock()
    		if p.err != nil {
    			// parent has already been canceled
    			child.cancel(false, p.err)
    		} else {
    			if p.children == nil {
    				p.children = make(map[canceler]struct{}) //创建指向children的set集合
    			}
    			p.children[child] = struct{}{}
    		}
    		p.mu.Unlock()
    	} else {
    		atomic.AddInt32(&goroutines, +1)
    		go func() {
    			select {
    			case <-parent.Done():
    				child.cancel(false, parent.Err())
    			case <-child.Done():
    			}
    		}()
    	}
    }
    ```

+ **valueCtx** 用于传递上下文共享信息

  ```go
  type valueCtx struct {
  	Context			//内嵌结构体，这个相当于指向父Context的指针
  	key, val any	//k-v, any是空接口（任意类型）
  }
  
  //查找共享的value，由源码可以看到k-v也可以覆盖
  func (c *valueCtx) Value(key any) any {
  	if c.key == key {	//先查当前节点的key value
  		return c.val
  	}
  	return value(c.Context, key)	//再查父节点key value
  }
  ```



## 参考资料

+ https://pkg.go.dev/context
+ [深度解密Go语言之context](https://zhuanlan.zhihu.com/p/68792989)