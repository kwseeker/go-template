package inbound

import (
	"context"
	"kwseeker.top/kwseeker/go-template/basic/net/socks/v2ray/common"
	"kwseeker.top/kwseeker/go-template/basic/net/socks/v2ray/features"
)

type Handler interface {
	common.Runnable
	Tag() string
}

type Manager interface {
	features.Feature

	GetHandler(ctx context.Context, tag string) (Handler, error)
	AddHandler(ctx context.Context, handler Handler) error
	RemoveHandler(ctx context.Context, tag string) error
}

func ManagerType() interface{} {
	return (*Manager)(nil)
}
