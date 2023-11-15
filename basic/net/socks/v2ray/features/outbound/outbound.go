package outbound

import (
	"context"
	"kwseeker.top/kwseeker/go-template/basic/net/socks/v2ray/common"
	"kwseeker.top/kwseeker/go-template/basic/net/socks/v2ray/features"
)

type Handler interface {
	common.Runnable
	Tag() string
	//Dispatch(ctx context.Context, link *transport.Link)
}

type Manager interface {
	features.Feature
	// GetHandler returns an outbound.Handler for the given tag.
	GetHandler(tag string) Handler
	// GetDefaultHandler returns the default outbound.Handler. It is usually the first outbound.Handler specified in the configuration.
	GetDefaultHandler() Handler
	// AddHandler adds a handler into this outbound.Manager.
	AddHandler(ctx context.Context, handler Handler) error
	// RemoveHandler removes a handler from outbound.Manager.
	RemoveHandler(ctx context.Context, tag string) error
}

func ManagerType() interface{} {
	return (*Manager)(nil)
}
