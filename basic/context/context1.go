package context

import (
	"context"
	"fmt"
)

/*
context 主要用来在 goroutine 之间传递上下文信息，包括：取消信号、超时时间、截止时间、k-v 等。
这里是“传递共享数据”的DEMO,比如requestId、traceId
*/

func process(ctx context.Context) {
	traceId, ok := ctx.Value("traceId").(string) //emptyCtx Value() 方法返回nil
	if ok {
		fmt.Printf("process over. trace_id=%s\n", traceId)
	} else {
		fmt.Printf("process over. no trace_id\n")
	}
}
