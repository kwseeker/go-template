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



## 对象池





## 应用场景



## 参考

+ [Golang 并发编程之同步原语](https://mp.weixin.qq.com/s?__biz=MzUzNTY5MzU2MA==&mid=2247484379&idx=1&sn=1a2abc6f639a34e62f3a5a0fcd774a71&chksm=fa80d24ccdf75b5a70d45168ad9e3a55dd258c1dd57147166a86062ee70d909ff1e5b0ba2bcb&token=183756123&lang=zh_CN#rd)