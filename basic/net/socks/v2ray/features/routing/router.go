package routing

import (
	"kwseeker.top/kwseeker/go-template/basic/net/socks/v2ray/common/errors"
	"kwseeker.top/kwseeker/go-template/basic/net/socks/v2ray/features"
)

type Router interface {
	features.Feature

	PickRoute(ctx Context) (Route, error)
}

type Route interface {
	Context
}

func RouterType() interface{} {
	return (*Router)(nil)
}

type DefaultRouter struct{}

// Type implements common.HasType.
func (DefaultRouter) Type() interface{} {
	return RouterType()
}

// PickRoute implements Router.
func (DefaultRouter) PickRoute(ctx Context) (Route, error) {
	return nil, errors.NewError("PickRoute method not implemented")
}

// Start implements common.Runnable.
func (DefaultRouter) Start() error {
	return nil
}

// Close implements common.Closable.
func (DefaultRouter) Close() error {
	return nil
}
