package expander

import "sync"

// Buffer is a simple thread-safe fixed-size buffer of runes.
// It keeps only the last size runes that were appended.
type Buffer struct {
	data []rune
	size int
	mu   sync.RWMutex
}

// NewBuffer creates a new buffer with the given maximum size.
// If size is non-positive, a default of 50 is used.
func NewBuffer(size int) *Buffer {
	if size <= 0 {
		size = 50
	}
	return &Buffer{
		data: make([]rune, 0, size),
		size: size,
	}
}

// Append adds a rune to the end of the buffer, discarding the
// oldest rune if the buffer is already full.
func (b *Buffer) Append(char rune) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.data) >= b.size {
		// Shift everything left by one and place the new char at the end.
		copy(b.data, b.data[1:])
		b.data[b.size-1] = char
		return
	}

	b.data = append(b.data, char)
}

// Remove removes the last rune from the buffer, if any.
func (b *Buffer) Remove() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.data) == 0 {
		return
	}
	b.data = b.data[:len(b.data)-1]
}

// Clear removes all runes from the buffer.
func (b *Buffer) Clear() {
	b.mu.Lock()
	b.data = b.data[:0]
	b.mu.Unlock()
}

// String returns the contents of the buffer as a string.
func (b *Buffer) String() string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return string(b.data)
}

// EndsWith reports whether the buffer currently ends with the given trigger string.
func (b *Buffer) EndsWith(trigger string) bool {
	if trigger == "" {
		return false
	}

	b.mu.RLock()
	defer b.mu.RUnlock()

	tr := []rune(trigger)
	if len(b.data) < len(tr) {
		return false
	}

	offset := len(b.data) - len(tr)
	for i, r := range tr {
		if b.data[offset+i] != r {
			return false
		}
	}
	return true
}