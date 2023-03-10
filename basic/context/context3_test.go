package context

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

// 应用中最常见的实现
// 这里模拟10后应用自动退出或收到关闭信号后退出
func TestAppShutdown(t *testing.T) {
	//ctx, cancel := context.WithCancel(context.Background())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	//<-ctx.Done()
	//fmt.Println("ctx done!")
	//<-s
	//fmt.Println("received sig!")

	select {
	case <-ctx.Done():
		cancel()
	case <-s:
		cancel()
	}
	fmt.Println("main out")
}
