package event

import (
	"sync"
)

type Handler func(payload any)

type Bus struct {
	mu       sync.RWMutex
	handlers map[string][]Handler
}

func NewBus() *Bus {
	return &Bus{
		handlers: make(map[string][]Handler),
	}
}

func (b *Bus) Subscribe(eventName string, handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventName] = append(b.handlers[eventName], handler)
}

func (b *Bus) Publish(eventName string, payload any) {
	b.mu.RLock()
	handlers := append([]Handler(nil), b.handlers[eventName]...)
	b.mu.RUnlock()
	for _, handler := range handlers {
		handler(payload)
	}
}

func (b *Bus) PublishAsync(eventName string, payload any) {
	b.mu.RLock()
	handlers := append([]Handler(nil), b.handlers[eventName]...)
	b.mu.RUnlock()
	for _, handler := range handlers {
		go handler(payload)
	}
}
