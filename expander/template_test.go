package expander

import (
	"strings"
	"testing"
)

func TestTemplateProcessorDate(t *testing.T) {
	tp := NewTemplateProcessor()
	result, offset := tp.Process("Today is {DATE}.")

	if !strings.HasPrefix(result, "Today is ") || !strings.HasSuffix(result, ".") {
		t.Fatalf("unexpected result: %q", result)
	}
	if offset != 0 {
		t.Fatalf("expected cursorOffset 0 when no {CURSOR}, got %d", offset)
	}
}

func TestTemplateProcessorCursor(t *testing.T) {
	tp := NewTemplateProcessor()
	result, offset := tp.Process("Hello{CURSOR}World")

	if result != "HelloWorld" {
		t.Fatalf("unexpected result: %q", result)
	}

	// Cursor should be placed before "World", i.e. offset equal to len("World")
	if offset != len("World") {
		t.Fatalf("expected cursorOffset %d, got %d", len("World"), offset)
	}
}

func TestTemplateProcessorCustomVar(t *testing.T) {
	tp := NewTemplateProcessor()
	tp.SetCustomVar("NAME", "Alice")
	result, _ := tp.Process("Hi {NAME}")

	if result != "Hi Alice" {
		t.Fatalf("unexpected result with custom var: %q", result)
	}
}