# sync

包结构（去除测试）：

```go
sync
├── atomic
│   ├── asm.s
│   ├── race.s
│   ├── type.go
│   ├── value.go
├── cond.go
├── map.go
├── mutex.go
├── once.go
├── pool.go
├── poolqueue.go
├── runtime2.go
├── runtime2_lockrank.go
├── runtime.go
├── rwmutex.go
├── waitgroup.go
```

sync 包提供了常见的并发编程同步原语,  包括 Mutex（互斥锁）、RWMutex（读写锁）、WaitGroup（）、Once 和 Cond， 以及拓展原语 ErrGroup、Semaphore和 SingleFlight 。

> 所谓**原语**，一般是指由若干条指令组成的程序段，用来实现某个特定功能，在执行过程中**不可被中断**。
>
> **Go 锁依旧是通过内存访问的方式进行并发同步的，与CSP是不同的并发同步方式。**

## 同步原语

### Mutex

```go
type Mutex struct {
	state int32			//互斥锁状态（4 bytes）(高29位：记录等待当前互斥锁的Goroutine个数，低3bit: mutexStarving、mutexWoken、mutexLocked)
    							  //
    sema  uint32	//控制锁状态的信号量(4 bytes)
}
```

只有3个公共方法

```go
Lock()						//阻塞等待获取锁
TryLock() bool		//尝试获取锁，返回获取结果
Unlock() 				//解锁

func (m *Mutex) Lock() {
	// Fast path: grab unlocked mutex.
	if atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked) {
		if race.Enabled {
			race.Acquire(unsafe.Pointer(m))
		}
		return
	}
	// Slow path (outlined so that the fast path can be inlined)
	m.lockSlow()
}
```

### RWMutex

### WaitGroup

等待一组并发操作完成。当不关心并发操作结果，或者有其他方法来收集它们的结果时使用。当这些条件不满足时应使用 channel 和 select 语句。

> WaitGroup 名字很贴切，和 Java 中的 CountLatch 和 Semphore 功能类似。

使用Channel 实现 WaitGroup 效果（但是可以通过 channel 传递结果）：

```go
done := make(chan bool)
for i := 0; i < 10; i++ {
    go func() {
        //do something ...
        done <- true
    }()
}
for i := 0; i < 10; i++ {
    <-done
}
```

### Cond

`sync.Cond` 基于互斥锁/读写锁，互斥锁 `sync.Mutex` 通常用来保护临界区和共享资源，条件变量 `sync.Cond` 用来协调想要访问共享资源的 goroutine。

功能都体现到方法上了（即 Go 的等待唤醒机制）。

```go
type Cond struct {
	noCopy noCopy

	// L is held while observing or changing the condition
	L Locker

	notify  notifyList
	checker copyChecker
}

func (c *Cond) Wait()
func (c *Cond) Signal()
func (c *Cond) Broadcast()
```

### Once

确保某个函数无论在单个协程或者多个协程中都只会被调用一次。



## Channel & Select 



## GOMAXPROCS 控制

此函数位于 runtime 包。`runtime.GOMAXPROCS`的作用是设置当前进程使用的最大cpu数，返回值为上一次调用成功的设置值。

Go1.5 之后默认使用CPU最大核心数，之前默认是使用一个核心。



## 并发同步选择：sync锁 vs CSP (channel)

参考《Go语言并发之道》图2-1。





## 应用场景



## 参考

+ [Golang 并发编程之同步原语](https://mp.weixin.qq.com/s?__biz=MzUzNTY5MzU2MA==&mid=2247484379&idx=1&sn=1a2abc6f639a34e62f3a5a0fcd774a71&chksm=fa80d24ccdf75b5a70d45168ad9e3a55dd258c1dd57147166a86062ee70d909ff1e5b0ba2bcb&token=183756123&lang=zh_CN#rd)