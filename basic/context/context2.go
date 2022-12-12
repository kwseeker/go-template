package context

import (
	"context"
	"fmt"
	"time"
)

/*
取消协程DEMO，Server启动时
*/

func Perform(ctx context.Context) {
	for {
		execBiz()
		sendResult()

		select {
		case <-ctx.Done():
			fmt.Printf(">>>>>>> been canceled, exit ...\n")
			return
		case <-time.After(time.Second):
			fmt.Printf("block 1 second ...\n")
		}
	}
}

func execBiz() {
	fmt.Printf("exec biz ...\n")
}

func sendResult() {
	fmt.Printf("send result to client ...\n")
}
