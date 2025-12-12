package expander

import (
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/atotto/clipboard"
)

// TemplateProcessor replaces variables in expansion templates and determines
// cursor placement using the {CURSOR} marker.
type TemplateProcessor struct {
	customVars map[string]string
	mu         sync.RWMutex
}

// NewTemplateProcessor creates a new template processor.
func NewTemplateProcessor() *TemplateProcessor {
	return &TemplateProcessor{
		customVars: make(map[string]string),
	}
}

// SetCustomVar sets a single custom variable key/value pair.
func (tp *TemplateProcessor) SetCustomVar(key, value string) {
	tp.mu.Lock()
	defer tp.mu.Unlock()

	if tp.customVars == nil {
		tp.customVars = make(map[string]string)
	}
	tp.customVars[key] = value
}

// SetCustomVars replaces the current custom variable map with the provided one.
func (tp *TemplateProcessor) SetCustomVars(vars map[string]string) {
	tp.mu.Lock()
	defer tp.mu.Unlock()

	if vars == nil {
		tp.customVars = make(map[string]string)
		return
	}

	cp := make(map[string]string, len(vars))
	for k, v := range vars {
		cp[k] = v
	}
	tp.customVars = cp
}

// GetAllVars returns a copy of all custom variables.
func (tp *TemplateProcessor) GetAllVars() map[string]string {
	tp.mu.RLock()
	defer tp.mu.RUnlock()

	cp := make(map[string]string, len(tp.customVars))
	for k, v := range tp.customVars {
		cp[k] = v
	}
	return cp
}

// Process applies all supported variables to the template and returns the final
// string along with the cursor offset from the end of the string. If no
// {CURSOR} marker is present, cursorOffset will be 0 (cursor at end).
func (tp *TemplateProcessor) Process(template string) (result string, cursorOffset int) {
	if template == "" {
		return "", 0
	}

	now := time.Now()
	dateStr := now.Format("2006-01-02")
	timeStr := now.Format("15:04:05")
	dateTimeStr := now.Format("2006-01-02 15:04:05")

	tp.mu.RLock()
	customCopy := make(map[string]string, len(tp.customVars))
	for k, v := range tp.customVars {
		customCopy[k] = v
	}
	tp.mu.RUnlock()

	var (
		builder           strings.Builder
		cursorPosInResult = -1 // rune index from start in final result
	)

	runes := []rune(template)
	for i := 0; i < len(runes); {
		if runes[i] == '{' {
			j := i + 1
			for j < len(runes) && runes[j] != '}' {
				j++
			}
			if j < len(runes) && runes[j] == '}' {
				token := string(runes[i+1 : j])
				upperToken := strings.ToUpper(token)

				// Handle built-in variables
				switch upperToken {
				case "DATE":
					builder.WriteString(dateStr)
				case "TIME":
					builder.WriteString(timeStr)
				case "DATETIME":
					builder.WriteString(dateTimeStr)
				case "CLIPBOARD":
					if text, err := clipboard.ReadAll(); err == nil {
						builder.WriteString(text)
					}
				case "CURSOR":
					// Mark position but do not output anything
					cursorPosInResult = utf8.RuneCountInString(builder.String())
				default:
					if val, ok := customCopy[upperToken]; ok {
						builder.WriteString(val)
					} else {
						// Unknown variable, keep literal text including braces.
						builder.WriteRune('{')
						builder.WriteString(token)
						builder.WriteRune('}')
					}
				}

				i = j + 1
				continue
			}
		}

		builder.WriteRune(runes[i])
		i++
	}

	result = builder.String()
	if cursorPosInResult < 0 {
		return result, 0
	}

	totalRunes := utf8.RuneCountInString(result)
	cursorOffset = totalRunes - cursorPosInResult
	if cursorOffset < 0 {
		cursorOffset = 0
	}
	return result, cursorOffset
}