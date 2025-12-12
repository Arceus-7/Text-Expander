package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Logger handles expansion logging and usage statistics.
type Logger struct {
	file   *os.File
	path   string
	mu     sync.Mutex
	stats  Statistics
	counts map[string]int
}

// Statistics summarises expansion usage.
type Statistics struct {
	TotalExpansions int
	TodayExpansions int
	MostUsedTrigger string
	LastExpansion   time.Time
}

const maxLogSize = 5 * 1024 * 1024 // 5MB

// NewLogger creates a new logger writing to the specified file path.
// If the logger cannot be created, nil is returned.
func NewLogger(path string) *Logger {
	if path == "" {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil
	}

	return &Logger{
		file:   f,
		path:   path,
		counts: make(map[string]int),
	}
}

// Close closes the underlying log file.
func (l *Logger) Close() {
	if l == nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file != nil {
		_ = l.file.Close()
		l.file = nil
	}
}

// LogExpansion records an expansion event, without logging the expanded text.
func (l *Logger) LogExpansion(trigger string) {
	if l == nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()

	// Update statistics.
	if !sameDay(l.stats.LastExpansion, now) {
		l.stats.TodayExpansions = 0
	}
	l.stats.TotalExpansions++
	l.stats.TodayExpansions++
	l.stats.LastExpansion = now

	l.counts[trigger]++
	if l.stats.MostUsedTrigger == "" || l.counts[trigger] > l.counts[l.stats.MostUsedTrigger] {
		l.stats.MostUsedTrigger = trigger
	}

	if l.file != nil {
		l.rotateIfNeeded()
		fmt.Fprintf(l.file, "%s\ttrigger=%s\n", now.Format(time.RFC3339), trigger)
	}
}

// LogError records an error message in the log.
func (l *Logger) LogError(err error) {
	if l == nil || err == nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file != nil {
		l.rotateIfNeeded()
		fmt.Fprintf(l.file, "%s\tERROR: %v\n", time.Now().Format(time.RFC3339), err)
	}
}

// GetStats returns a snapshot of the current statistics.
func (l *Logger) GetStats() Statistics {
	if l == nil {
		return Statistics{}
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	return l.stats
}

func (l *Logger) rotateIfNeeded() {
	if l.file == nil {
		return
	}

	info, err := l.file.Stat()
	if err != nil {
		return
	}
	if info.Size() < maxLogSize {
		return
	}

	_ = l.file.Close()
	backupName := fmt.Sprintf("%s.%d", l.path, time.Now().Unix())
	_ = os.Rename(l.path, backupName)

	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		l.file = nil
		return
	}
	l.file = f
}

func sameDay(a, b time.Time) bool {
	if a.IsZero() || b.IsZero() {
		return false
	}
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}