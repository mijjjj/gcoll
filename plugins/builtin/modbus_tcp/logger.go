package modbustcp

import (
	"sync"
	"time"
)

// DebugLog 描述插件调试日志。
type DebugLog struct {
	Time     time.Time         `json:"time"`
	Level    string            `json:"level"`
	DeviceID string            `json:"deviceId"`
	TaskID   string            `json:"taskId"`
	PointID  string            `json:"pointId"`
	Message  string            `json:"message"`
	TraceID  string            `json:"traceId"`
	Fields   map[string]string `json:"fields"`
}

// DebugBuffer 保存最近的调试日志。
type DebugBuffer struct {
	mu      sync.Mutex
	limit   int
	enabled bool
	items   []DebugLog
}

// NewDebugBuffer 创建调试日志缓冲区。
func NewDebugBuffer(limit int, enabled bool) *DebugBuffer {
	if limit <= 0 {
		limit = defaultDebugLogLimit
	}
	return &DebugBuffer{limit: limit, enabled: enabled}
}

// SetEnabled 设置调试模式。
func (b *DebugBuffer) SetEnabled(enabled bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.enabled = enabled
}

// Append 追加调试日志。
func (b *DebugBuffer) Append(log DebugLog) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if !b.enabled {
		return
	}
	if log.Time.IsZero() {
		log.Time = time.Now()
	}
	b.items = append(b.items, log)
	if len(b.items) > b.limit {
		b.items = b.items[len(b.items)-b.limit:]
	}
}

// Snapshot 返回最近日志快照。
func (b *DebugBuffer) Snapshot() []DebugLog {
	b.mu.Lock()
	defer b.mu.Unlock()
	return append([]DebugLog(nil), b.items...)
}
