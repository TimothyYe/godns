package lib

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// LogEntry represents a single log entry.
type LogEntry struct {
	Timestamp time.Time     `json:"timestamp"`
	Level     string        `json:"level"`
	Message   string        `json:"message"`
	Fields    logrus.Fields `json:"fields,omitempty"`
}

// LogBuffer is a thread-safe circular buffer for log entries.
type LogBuffer struct {
	entries []LogEntry
	size    int
	index   int
	mutex   sync.RWMutex
	full    bool
}

// NewLogBuffer creates a new log buffer with specified size.
func NewLogBuffer(size int) *LogBuffer {
	return &LogBuffer{
		entries: make([]LogEntry, size),
		size:    size,
		index:   0,
		full:    false,
	}
}

// Add adds a new log entry to the buffer.
func (lb *LogBuffer) Add(entry LogEntry) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	lb.entries[lb.index] = entry
	lb.index = (lb.index + 1) % lb.size

	if lb.index == 0 {
		lb.full = true
	}
}

// GetAll returns all log entries in chronological order.
func (lb *LogBuffer) GetAll() []LogEntry {
	lb.mutex.RLock()
	defer lb.mutex.RUnlock()

	if !lb.full && lb.index == 0 {
		return []LogEntry{}
	}

	var result []LogEntry

	if lb.full {
		// Buffer is full, start from current index (oldest entry)
		result = make([]LogEntry, 0, lb.size)
		for i := 0; i < lb.size; i++ {
			idx := (lb.index + i) % lb.size
			result = append(result, lb.entries[idx])
		}
	} else {
		// Buffer is not full, return entries from 0 to index
		result = make([]LogEntry, lb.index)
		copy(result, lb.entries[:lb.index])
	}

	return result
}

// GetRecent returns the most recent n log entries.
func (lb *LogBuffer) GetRecent(n int) []LogEntry {
	all := lb.GetAll()
	if len(all) <= n {
		return all
	}
	return all[len(all)-n:]
}

// Clear clears all log entries.
func (lb *LogBuffer) Clear() {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	lb.index = 0
	lb.full = false
	lb.entries = make([]LogEntry, lb.size)
}

// Global log buffer instance.
var globalLogBuffer *LogBuffer

// InitLogBuffer initializes the global log buffer.
func InitLogBuffer(size int) {
	globalLogBuffer = NewLogBuffer(size)
}

// GetLogBuffer returns the global log buffer instance.
func GetLogBuffer() *LogBuffer {
	if globalLogBuffer == nil {
		InitLogBuffer(1000) // Default size
	}
	return globalLogBuffer
}

// LogHook is a logrus hook that captures logs in the buffer.
type LogHook struct {
	buffer *LogBuffer
}

// NewLogHook creates a new log hook.
func NewLogHook(buffer *LogBuffer) *LogHook {
	return &LogHook{
		buffer: buffer,
	}
}

// Fire is called when a log event is fired.
func (hook *LogHook) Fire(entry *logrus.Entry) error {
	logEntry := LogEntry{
		Timestamp: entry.Time,
		Level:     entry.Level.String(),
		Message:   entry.Message,
		Fields:    make(logrus.Fields),
	}

	// Copy fields
	for k, v := range entry.Data {
		logEntry.Fields[k] = v
	}

	hook.buffer.Add(logEntry)
	return nil
}

// Levels returns the available logging levels.
func (hook *LogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
