package expander

import "testing"

func TestBufferEndsWith(t *testing.T) {
	b := NewBuffer(50)
	for _, r := range []rune(";email") {
		b.Append(r)
	}

	if !b.EndsWith(";email") {
		t.Error("Buffer should end with ;email")
	}
}

func TestBufferAppendAndRemove(t *testing.T) {
	b := NewBuffer(3)
	b.Append('a')
	b.Append('b')
	b.Append('c')

	if got := b.String(); got != "abc" {
		t.Fatalf("expected abc, got %q", got)
	}

	b.Append('d')
	if got := b.String(); got != "bcd" {
		t.Fatalf("expected bcd after overflow, got %q", got)
	}

	b.Remove()
	if got := b.String(); got != "bc" {
		t.Fatalf("expected bc after remove, got %q", got)
	}

	b.Clear()
	if got := b.String(); got != "" {
		t.Fatalf("expected empty string after Clear, got %q", got)
	}
}