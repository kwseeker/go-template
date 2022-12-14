# Pprof

Go 引入了 Pprof 作为分析性能、分析数据的工具，Pprof是一个用于概要数据可视化和分析的工具。Pprof读取概要文件中的分析样本集合。Proto格式和生成报告，以可视化和帮助分析数据。它可以生成文本和图形报告(通过使用点可视化包)。

用途：

- **CPU Profiling**：CPU 分析，按照一定的频率采集所监听的应用程序 CPU（含寄存器）的使用情况，可确定应用程序在主动消耗 CPU 周期时花费时间的位置。
- **Memory Profiling**：内存分析，在应用程序进行堆分配时记录堆栈跟踪，用于监视当前和历史内存使用情况，以及检查内存泄漏。
- **Block Profiling**：阻塞分析，记录Goroutine阻塞等待同步（包括定时器通道）的位置，默认不开启，需要调用`runtime.SetBlockProfileRate`进行设置。
- **Mutex Profiling**：互斥锁分析，报告互斥锁的竞争情况，默认不开启，需要调用`runtime.SetMutexProfileFraction`进行设置。
- **Goroutine Profiling**：Goroutine 分析，可以对当前应用程序正在运行的 Goroutine 进行堆栈跟踪和分析。

Go Pprof 3种采样方式：

- **runtime/pprof**：采集程序（非 Server）的指定区块的运行数据进行分析。

- **net/http/pprof**：基于HTTP Server运行，并且可以采集运行时数据进行分析。

  ```go
   import _ "net/http/pprof"
   //代码中启动HTTP server
   _ = http.ListenAndServe("0.0.0.0:6060", nil)
  ```

- **go test**：通过运行测试用例，并指定所需标识来进行采集。

Go Pprof 3种使用模式：

- Report generation：报告生成。
- Interactive terminal use：交互式终端使用。
- Web interface：Web 界面。



## 案例

+ **SDK mutex_test.go中分析互斥锁**

  ```go
  func TestMutex(t *testing.T) {
      //开启互斥锁分析，报告互斥锁的竞争情况,rate控制采样率(是分母),这里1即采样率1/1，即100%的采样率 
  	if n := runtime.SetMutexProfileFraction(1); n != 0 {
  		t.Logf("got mutexrate %d expected 0", n)
  	}
  	defer runtime.SetMutexProfileFraction(0)
      runtime.SetBlockProfileRate(1)
  	defer runtime.SetBlockProfileRate(0)
      ......
}
  ```
  
  交互式终端分析profile文件：
  
  ```
  $ go test -v -mutexprofile=mutex.prof -blockprofile=block.prof
  $ go tool pprof block.prof                                    
  File: sync.test
  Type: delay
  Time: Dec 13, 2022 at 12:41pm (CST)
  Entering interactive mode (type "help" for commands, "o" for options)
  (pprof) top
  Showing nodes accounting for 148.32ms, 99.95% of 148.40ms total
  Dropped 1 node (cum <= 0.74ms)
        flat  flat%   sum%        cum   cum%
    145.55ms 98.08% 98.08%   145.55ms 98.08%  sync.(*Mutex).Lock (inline)
      2.77ms  1.87% 99.95%     2.77ms  1.87%  runtime.chanrecv1
           0     0% 99.95%   145.63ms 98.13%  kwseeker.top/kwseeker/go-template/basic/sync.HammerMutex
           0     0% 99.95%     2.77ms  1.87%  kwseeker.top/kwseeker/go-template/basic/sync.TestMutex
           0     0% 99.95%     2.77ms  1.87%  testing.tRunner
  (pprof) list HammerMutex
  Total: 148.40ms
  ROUTINE ======================== kwseeker.top/kwseeker/go-template/basic/sync.HammerMutex in /home/lee/go/src/kwseeker.top/kwseeker/go-template/basic/sync/mutex_test.go
           0   145.63ms (flat, cum) 98.13% of Total
           .          .     13:                           *sum++
           .          .     14:                           m.Unlock()
           .          .     15:                   }
           .          .     16:                   continue
           .          .     17:           }
           .   145.55ms     18:           m.Lock() //阻塞等待获取锁
           .          .     19:           *sum++
           .          .     20:           m.Unlock()
           .          .     21:   }
           .    76.58us     22:   cdone <- true
           .          .     23:}
           .          .     24:
           .          .     25:/*
           .          .     26:go test -v -mutexprofile=mutex.prof -blockprofile=block.prof
           .          .     27:go tool pprof mutex.prof
  (pprof) web
  展示在浏览器
  ```



## 参考

+ [Go 程序崩了？煎鱼教你用 PProf 工具来救火！](https://zhuanlan.zhihu.com/p/409592769)