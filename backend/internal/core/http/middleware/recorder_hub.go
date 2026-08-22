package middleware

import "sync"

type RecorderHub struct {
	mu sync.RWMutex
	fn RequestRecorder
}

func NewRecorderHub() *RecorderHub {
	return &RecorderHub{}
}

func (h *RecorderHub) Set(fn RequestRecorder) {
	if h == nil {
		return
	}
	h.mu.Lock()
	h.fn = fn
	h.mu.Unlock()
}

func (h *RecorderHub) Record(entry RequestLogEntry) {
	if h == nil {
		return
	}
	h.mu.RLock()
	fn := h.fn
	h.mu.RUnlock()
	if fn != nil {
		fn(entry)
	}
}
