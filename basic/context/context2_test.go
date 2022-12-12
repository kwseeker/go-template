package context

import (
	"context"
	"fmt"
	"math"
	"testing"
	"time"
)

func TestCtx2(t *testing.T) {
	//模拟一个请求
	go func() {
		//ctx, cancel := context.WithTimeout(context.Background(), time.Minute) //超时自动调用cancel
		ctx, cancel := context.WithCancel(context.Background())
		go Perform(ctx)

		// ……
		time.Sleep(time.Second * 5)
		//time.Sleep(time.Second * 65)

		// 任务结束，主动调用cancel 函数
		fmt.Printf(">>>>>>> call cancel ...\n")
		cancel()
	}()

	time.Sleep(time.Second * math.MaxInt32)
}
