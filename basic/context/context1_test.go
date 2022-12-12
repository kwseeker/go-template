package context

import (
	"context"
	"testing"
)

func TestCtx1(t *testing.T) {
	//Background()函数创建空的根上下文, 本质是个emptyCtx实例, 方法基本都是空的
	ctx := context.Background()
	process(ctx)

	//WithValue()函数创建valueCtx &valueCtx{parent, key, val} valueCtx内嵌结构体Context, 正好变成向着根方向的指针
	ctx = context.WithValue(ctx, "traceId", "10001")
	process(ctx)
}
