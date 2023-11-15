package inbound

import (
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	errors2 "kwseeker.top/kwseeker/go-template/basic/net/socks/v2ray/common/errors"
	"kwseeker.top/kwseeker/go-template/basic/net/socks/v2ray/common/serial"
	"kwseeker.top/kwseeker/go-template/basic/net/socks/v2ray/features/inbound"
	"sync"
)

// Manager 管理所有 inbound Handler
// 实现 features/inbound/inbound.Manager 接口
type Manager struct {
	ctx             context.Context
	access          sync.RWMutex
	untaggedHandler []inbound.Handler
	taggedHandlers  map[string]inbound.Handler
	running         bool
}

func New(ctx context.Context) (*Manager, error) {
	m := &Manager{
		ctx:            ctx,
		taggedHandlers: make(map[string]inbound.Handler),
	}
	return m, nil
}

func (m *Manager) Type() interface{} {
	return inbound.ManagerType()
}

// Start implements common.Runnable.
func (m *Manager) Start() error {
	m.access.Lock()
	defer m.access.Unlock()

	m.running = true

	for _, handler := range m.taggedHandlers {
		if err := handler.Start(); err != nil {
			return err
		}
	}

	for _, handler := range m.untaggedHandler {
		if err := handler.Start(); err != nil {
			return err
		}
	}
	return nil
}

// Close implements common.Closable.
func (m *Manager) Close() error {
	m.access.Lock()
	defer m.access.Unlock()

	m.running = false

	var errs []interface{}
	for _, handler := range m.taggedHandlers {
		if err := handler.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	for _, handler := range m.untaggedHandler {
		if err := handler.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.New(fmt.Sprint("failed to close all handlers", serial.Concat(errs)))
	}

	return nil
}

func (m *Manager) GetHandler(ctx context.Context, tag string) (inbound.Handler, error) {
	m.access.RLock()
	defer m.access.RUnlock()

	handler, found := m.taggedHandlers[tag]
	if !found {
		return nil, errors.New(fmt.Sprint("handler not found: ", tag))
	}
	return handler, nil
}

func (m *Manager) AddHandler(ctx context.Context, handler inbound.Handler) error {
	m.access.Lock()
	defer m.access.Unlock()

	tag := handler.Tag()
	if len(tag) > 0 {
		m.taggedHandlers[tag] = handler
	} else {
		m.untaggedHandler = append(m.untaggedHandler, handler)
	}

	if m.running {
		return handler.Start()
	}

	return nil
}

func (m *Manager) RemoveHandler(ctx context.Context, tag string) error {
	if tag == "" {
		return errors.New("empty tag when remove inbound handler")
	}

	m.access.Lock()
	defer m.access.Unlock()

	if handler, found := m.taggedHandlers[tag]; found {
		if err := handler.Close(); err != nil {
			log.Error("failed to close handler ", tag)
		}
		delete(m.taggedHandlers, tag)
		return nil
	}

	return errors2.NewError("not found inbound handler, tag=", tag)
}
