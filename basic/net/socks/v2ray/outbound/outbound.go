package outbound

import (
	"context"
	"kwseeker.top/kwseeker/go-template/basic/net/socks/v2ray/common/errors"
	"kwseeker.top/kwseeker/go-template/basic/net/socks/v2ray/features/outbound"
	"log"
	"sync"
)

// Manager 管理所有 outbound Handler
// 实现 features/inbound/outbound.Manager 接口
type Manager struct {
	access           sync.RWMutex
	defaultHandler   outbound.Handler
	taggedHandler    map[string]outbound.Handler
	untaggedHandlers []outbound.Handler
	running          bool
}

func New() (*Manager, error) {
	m := &Manager{
		taggedHandler:    make(map[string]outbound.Handler),
		untaggedHandlers: make([]outbound.Handler, 0),
	}
	return m, nil
}

func (m *Manager) Type() interface{} {
	return outbound.ManagerType()
}

// Start implements core.Feature
func (m *Manager) Start() error {
	m.access.Lock()
	defer m.access.Unlock()

	m.running = true

	for _, h := range m.taggedHandler {
		if err := h.Start(); err != nil {
			return err
		}
	}

	for _, h := range m.untaggedHandlers {
		if err := h.Start(); err != nil {
			return err
		}
	}

	return nil
}

// Close implements core.Feature
func (m *Manager) Close() error {
	m.access.Lock()
	defer m.access.Unlock()

	m.running = false

	var errs []error
	for _, h := range m.taggedHandler {
		errs = append(errs, h.Close())
	}

	for _, h := range m.untaggedHandlers {
		errs = append(errs, h.Close())
	}

	return errors.Combine(errs...)
}

// GetDefaultHandler implements outbound.Manager.
func (m *Manager) GetDefaultHandler() outbound.Handler {
	m.access.RLock()
	defer m.access.RUnlock()

	return m.defaultHandler
}

// GetHandler implements outbound.Manager.
func (m *Manager) GetHandler(tag string) outbound.Handler {
	m.access.RLock()
	defer m.access.RUnlock()
	if handler, found := m.taggedHandler[tag]; found {
		return handler
	}
	return nil
}

// AddHandler implements outbound.Manager.
func (m *Manager) AddHandler(ctx context.Context, handler outbound.Handler) error {
	m.access.Lock()
	defer m.access.Unlock()
	tag := handler.Tag()

	if m.defaultHandler == nil ||
		(len(tag) > 0 && tag == m.defaultHandler.Tag()) {
		m.defaultHandler = handler
	}

	if len(tag) > 0 {
		if oldHandler, found := m.taggedHandler[tag]; found {
			log.Println("will replace the existed outbound with the tag: ", tag)
			_ = oldHandler.Close()
		}
		m.taggedHandler[tag] = handler
	} else {
		m.untaggedHandlers = append(m.untaggedHandlers, handler)
	}

	if m.running {
		return handler.Start()
	}

	return nil
}

// RemoveHandler implements outbound.Manager.
func (m *Manager) RemoveHandler(ctx context.Context, tag string) error {
	if tag == "" {
		return errors.NewError("empty tag when remove outbound handler")
	}
	m.access.Lock()
	defer m.access.Unlock()

	delete(m.taggedHandler, tag)
	if m.defaultHandler != nil && m.defaultHandler.Tag() == tag {
		m.defaultHandler = nil
	}

	return nil
}
