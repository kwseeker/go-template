package runtime

import (
	"os"
	"runtime/trace"
)

// TestTrace go run trace_example.go 2> trace.out	//将标准错误信息写入trace.out
func TestTrace() {
	//启动trace,将trace信息缓存并写入标准错误
	err := trace.Start(os.Stderr)
	if err != nil {
		return
	}
	defer trace.Stop()

	ch := make(chan string)
	go func() {
		ch <- "main out"
	}()

	<-ch
}
