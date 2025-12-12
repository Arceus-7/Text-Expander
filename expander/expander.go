package expander

import (
	"log"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/go-vgo/robotgo"

	"github.com/yourusername/text-expander/config"
	"github.com/yourusername/text-expander/utils"
)

// Expansion is an alias to the config-level Expansion type for convenience.
type Expansion = config.Expansion

// Expander ties together buffer management, keyboard hooks, template
// processing, and configuration to implement text expansion.
type Expander struct {
	expansions map[string]Expansion
	buffer     *Buffer
	keyboard   Keyboard
	config     *config.Config
	template   *TemplateProcessor
	logger     *utils.Logger

	mu          sync.RWMutex
	running     bool
	inExpansion bool
}

// NewExpander constructs a new Expander for the given configuration.
func NewExpander(cfg *config.Config) *Expander {
	return NewExpanderWithKeyboard(cfg, NewKeyboardHook())
}

// NewExpanderWithKeyboard allows injecting a custom keyboard implementation;
// this is primarily useful for testing.
func NewExpanderWithKeyboard(cfg *config.Config, kb Keyboard) *Expander {
	e := &Expander{
		expansions: make(map[string]Expansion),
		buffer:     NewBuffer(50),
		keyboard:   kb,
		config:     cfg,
		template:   NewTemplateProcessor(),
	}

	e.reloadFromConfigLocked()

	// Wire keyboard callback.
	if hook, ok := kb.(*KeyboardHook); ok {
		hook.SetOnKeyPress(e.OnKeyPress)
	}

	// Watch config file for changes and hot-reload.
	if cfg != nil {
		_ = cfg.Watch(func() {
			e.ReloadConfig()
		})
	}

	return e
}

// SetLogger attaches a logger to the expander for usage statistics.
func (e *Expander) SetLogger(l *utils.Logger) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.logger = l
}

// Start begins monitoring keyboard events.
func (e *Expander) Start() error {
	e.mu.Lock()
	if e.running {
		e.mu.Unlock()
		return nil
	}
	e.running = true
	kb := e.keyboard
	e.mu.Unlock()

	if kb == nil {
		return nil
	}
	return kb.Start()
}

// Stop stops monitoring keyboard events.
func (e *Expander) Stop() {
	e.mu.Lock()
	if !e.running {
		e.mu.Unlock()
		return
	}
	e.running = false
	kb := e.keyboard
	e.mu.Unlock()

	if kb != nil {
		kb.Stop()
	}
}

// OnKeyPress is invoked by the keyboard hook for each key press.
func (e *Expander) OnKeyPress(key string) {
	// Avoid processing keys that are produced by our own simulated typing.
	e.mu.RLock()
	if e.inExpansion {
		e.mu.RUnlock()
		return
	}
	cfg := e.config
	e.mu.RUnlock()

	if cfg == nil {
		log.Printf("[DEBUG] OnKeyPress: config is nil, ignoring key: %s", key)
		return
	}
	settings := cfg.GetSettings()

	// Debug: log key presses (limit to avoid spam)
	if key != KeySpace && key != KeyEnter && key != KeyTab && key != KeyBackspace {
		log.Printf("[DEBUG] OnKeyPress: received key: %q, buffer before: %q", key, e.buffer.String())
	}

	switch key {
	case KeyBackspace:
		e.buffer.Remove()
		log.Printf("[DEBUG] OnKeyPress: BACKSPACE, buffer after: %q", e.buffer.String())
	case KeySpace:
		bufBefore := e.buffer.String()
		if settings.Enabled && settings.TriggerOnSpace {
			log.Printf("[DEBUG] OnKeyPress: SPACE, checking expansion, buffer: %q", bufBefore)
			e.CheckAndExpand()
		} else {
			log.Printf("[DEBUG] OnKeyPress: SPACE, expansion disabled (enabled=%v, triggerOnSpace=%v)", settings.Enabled, settings.TriggerOnSpace)
		}
		e.buffer.Append(' ')
	case KeyEnter:
		bufBefore := e.buffer.String()
		if settings.Enabled && settings.TriggerOnEnter {
			log.Printf("[DEBUG] OnKeyPress: ENTER, checking expansion, buffer: %q", bufBefore)
			e.CheckAndExpand()
		} else {
			log.Printf("[DEBUG] OnKeyPress: ENTER, expansion disabled (enabled=%v, triggerOnEnter=%v)", settings.Enabled, settings.TriggerOnEnter)
		}
		e.buffer.Append('\n')
	case KeyTab:
		bufBefore := e.buffer.String()
		if settings.Enabled && settings.TriggerOnTab {
			log.Printf("[DEBUG] OnKeyPress: TAB, checking expansion, buffer: %q", bufBefore)
			e.CheckAndExpand()
		} else {
			log.Printf("[DEBUG] OnKeyPress: TAB, expansion disabled (enabled=%v, triggerOnTab=%v)", settings.Enabled, settings.TriggerOnTab)
		}
		e.buffer.Append('\t')
	default:
		runes := []rune(key)
		if len(runes) == 1 {
			e.buffer.Append(runes[0])
		}
	}
}

// CheckAndExpand inspects the input buffer for any matching trigger and, if
// found, performs the expansion.
func (e *Expander) CheckAndExpand() {
	bufferContent := e.buffer.String()
	log.Printf("[DEBUG] CheckAndExpand: buffer content: %q", bufferContent)

	if bufferContent == "" {
		log.Printf("[DEBUG] CheckAndExpand: buffer is empty, returning")
		return
	}

	e.mu.RLock()
	expansions := e.expansions
	e.mu.RUnlock()

	log.Printf("[DEBUG] CheckAndExpand: checking %d expansions", len(expansions))

	if len(expansions) == 0 {
		log.Printf("[DEBUG] CheckAndExpand: no expansions configured, returning")
		return
	}

	exp, ok := matchExpansion(bufferContent, expansions)
	if !ok {
		log.Printf("[DEBUG] CheckAndExpand: no matching expansion found for: %q", bufferContent)
		return
	}

	log.Printf("[DEBUG] CheckAndExpand: found match! trigger: %q, replacement: %q", exp.Trigger, exp.Replacement)
	e.PerformExpansion(exp.Trigger, exp.Replacement)
}

// PerformExpansion executes the delete-and-type sequence for a given trigger
// and expansion template.
func (e *Expander) PerformExpansion(trigger, expansion string) {
	if trigger == "" || expansion == "" {
		return
	}

	if !utils.ShouldAllowExpansion() {
		return
	}

	e.mu.RLock()
	tp := e.template
	logger := e.logger
	cfg := e.config
	e.mu.RUnlock()

	if tp == nil {
		return
	}

	text, cursorOffset := tp.Process(expansion)
	if text == "" {
		return
	}

	// Signal that we're in the middle of an expansion so we can ignore
	// synthetic key events from robotgo.
	e.mu.Lock()
	e.inExpansion = true
	e.mu.Unlock()

	defer func() {
		e.mu.Lock()
		e.inExpansion = false
		e.mu.Unlock()
	}()

	triggerLen := utf8.RuneCountInString(trigger)

	// Delete the trigger.
	if triggerLen > 0 && e.keyboard != nil {
		e.keyboard.SimulateBackspace(triggerLen)
		for i := 0; i < triggerLen; i++ {
			e.buffer.Remove()
		}
	}

	// Type the replacement.
	if e.keyboard != nil {
		e.keyboard.SimulateTyping(text)
		for _, r := range []rune(text) {
			e.buffer.Append(r)
		}
	}

	// Move cursor to requested position using left-arrow taps.
	for i := 0; i < cursorOffset; i++ {
		_ = robotgo.KeyTap(robotgo.Left)
		time.Sleep(2 * time.Millisecond)
	}

	// Log usage.
	if logger != nil && cfg != nil && cfg.GetSettings().LogExpansions {
		logger.LogExpansion(trigger)
	}
}

// ReloadConfig reloads configuration from disk when the config file changes.
func (e *Expander) ReloadConfig() {
	e.mu.RLock()
	cfg := e.config
	e.mu.RUnlock()

	if cfg == nil {
		return
	}

	path := cfg.ConfigPath()
	if path == "" {
		return
	}

	newCfg, err := config.LoadConfig(path)
	if err != nil {
		if l := e.logger; l != nil {
			l.LogError(err)
		}
		return
	}

	e.mu.Lock()
	e.config = newCfg
	e.reloadFromConfigLocked()
	e.mu.Unlock()
}

// reloadFromConfigLocked refreshes internal state from the current config.
// e.mu must be held by the caller.
func (e *Expander) reloadFromConfigLocked() {
	if e.config == nil {
		return
	}

	exps := e.config.GetExpansions()
	m := make(map[string]Expansion, len(exps))
	for _, exp := range exps {
		m[exp.Trigger] = exp
	}
	e.expansions = m

	if e.template != nil {
		e.template.SetCustomVars(e.config.GetCustomVars())
	}
}

// matchExpansion finds the best matching expansion for the given buffer
// content, preferring the longest trigger when multiple triggers match.
func matchExpansion(bufferContent string, expansions map[string]Expansion) (Expansion, bool) {
	if len(expansions) == 0 {
		return Expansion{}, false
	}

	lowerBuf := strings.ToLower(bufferContent)

	var selected Expansion
	found := false
	selectedLen := 0

	for _, exp := range expansions {
		if exp.Trigger == "" {
			continue
		}

		triggerLen := utf8.RuneCountInString(exp.Trigger)
		if triggerLen == 0 || triggerLen > utf8.RuneCountInString(bufferContent) {
			continue
		}

		if exp.CaseSensitive {
			if !strings.HasSuffix(bufferContent, exp.Trigger) {
				continue
			}
		} else {
			if !strings.HasSuffix(lowerBuf, strings.ToLower(exp.Trigger)) {
				continue
			}
		}

		if !found || triggerLen > selectedLen {
			selected = exp
			selectedLen = triggerLen
			found = true
		}
	}

	return selected, found
}
