# os.signal

信号的订阅 signal.Notify() 是通过 channel 完成的，每个`os.Signal` channel 都会收听自己的事件集。

![](../../img/signal原理.png)