package expander

import (
	"log"
	"sync"
	"time"

	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

// Logical key values passed to the expander.
const (
	KeyBackspace = "BACKSPACE"
	KeyEnter     = "ENTER"
	KeySpace     = "SPACE"
	KeyTab       = "TAB"
)

// Keyboard defines the subset of keyboard operations that the expander needs.
// It is implemented by KeyboardHook and can be replaced in tests.
type Keyboard interface {
	Start() error
	Stop()
	SimulateBackspace(count int)
	SimulateTyping(text string)
}

// KeyboardHook listens for global keyboard events using gohook and provides
// helpers to simulate key presses using robotgo.
type KeyboardHook struct {
	onKeyPress func(key string)

	mu      sync.Mutex
	running bool
}

// NewKeyboardHook creates a new keyboard hook instance.
func NewKeyboardHook() *KeyboardHook {
	return &KeyboardHook{}
}

// SetOnKeyPress sets the callback invoked on each key press.
// The callback is called from a background goroutine.
func (k *KeyboardHook) SetOnKeyPress(cb func(key string)) {
	k.mu.Lock()
	k.onKeyPress = cb
	k.mu.Unlock()
}

// Start begins listening for global keyboard events.
func (k *KeyboardHook) Start() error {
	k.mu.Lock()
	if k.running {
		k.mu.Unlock()
		return nil
	}
	k.running = true
	k.mu.Unlock()

	go func() {
		log.Printf("[DEBUG] KeyboardHook: starting keyboard hook...")
		evChan := hook.Start()
		defer hook.End()
		log.Printf("[DEBUG] KeyboardHook: keyboard hook started, waiting for events...")

		eventCount := 0
		for ev := range evChan {
			eventCount++
			// We only care about key down events.
			if ev.Kind != hook.KeyDown {
				continue
			}

			keyStr := hook.RawcodetoKeychar(ev.Rawcode)
			key := translateEvent(ev)

			// Always log space, enter, tab, backspace to debug translation
			if key == KeySpace || key == KeyEnter || key == KeyTab || key == KeyBackspace ||
				key == " " || key == "\r" || key == "\n" || key == "\t" || key == "" {
				log.Printf("[DEBUG] KeyboardHook: event #%d: keyStr=%q, keychar=%d, translated=%q", eventCount, keyStr, ev.Keychar, key)
			}

			if key == "" {
				continue
			}

			// Log first few events to verify hook is working
			if eventCount <= 10 {
				log.Printf("[DEBUG] KeyboardHook: received key event #%d: %q (rawcode=%d, keychar=%d, keyStr=%q)", eventCount, key, ev.Rawcode, ev.Keychar, keyStr)
			}

			k.mu.Lock()
			cb := k.onKeyPress
			k.mu.Unlock()

			if cb != nil {
				cb(key)
			} else {
				log.Printf("[DEBUG] KeyboardHook: WARNING - callback is nil!")
			}
		}
		log.Printf("[DEBUG] KeyboardHook: event channel closed, hook ending")
	}()

	return nil
}

// Stop stops listening for keyboard events.
func (k *KeyboardHook) Stop() {
	k.mu.Lock()
	defer k.mu.Unlock()

	if !k.running {
		return
	}
	hook.End()
	k.running = false
}

// SimulateBackspace sends the given number of backspace key taps.
func (k *KeyboardHook) SimulateBackspace(count int) {
	if count <= 0 {
		return
	}

	for i := 0; i < count; i++ {
		_ = robotgo.KeyTap(robotgo.Backspace)
		time.Sleep(2 * time.Millisecond)
	}
}

// SimulateTyping types the given text with small delays between keystrokes.
func (k *KeyboardHook) SimulateTyping(text string) {
	if text == "" {
		return
	}
	// TypeDelay introduces a per-key delay in milliseconds.
	robotgo.TypeDelay(text, 2)
}

// translateEvent maps a gohook Event to a logical key representation used by
// the expander.
func translateEvent(ev hook.Event) string {
	keyStr := hook.RawcodetoKeychar(ev.Rawcode)
	char := string(ev.Keychar)

	// Debug logging for special keys - log ALL potential special key events
	if ev.Keychar == 32 || ev.Keychar == 13 || ev.Keychar == 10 || ev.Keychar == 9 || ev.Keychar == 8 ||
		keyStr == " " || keyStr == "\r" || keyStr == "\n" || keyStr == "\t" || keyStr == "space" || keyStr == "enter" || keyStr == "tab" || keyStr == "backspace" {
		log.Printf("[DEBUG] translateEvent: keyStr=%q, keychar=%d, char=%q", keyStr, ev.Keychar, char)
	}

	// First check keyStr for named keys (this is more reliable)
	switch keyStr {
	case "space":
		return KeySpace
	case "enter":
		return KeyEnter
	case "tab":
		return KeyTab
	case "backspace":
		return KeyBackspace
	case "":
		// If keyStr is empty, check keychar for special keys
		if ev.Keychar == 0 || ev.Keychar == hook.CharUndefined {
			return ""
		}
		// Convert keychar to string first to check what it is
		char := string(ev.Keychar)

		// Check for special keys by their character representation
		// Space character
		if char == " " || ev.Keychar == 32 {
			return KeySpace
		}
		// Enter (carriage return or newline)
		if char == "\r" || char == "\n" || ev.Keychar == 13 || ev.Keychar == 10 {
			return KeyEnter
		}
		// Tab
		if char == "\t" || ev.Keychar == 9 {
			return KeyTab
		}
		// Backspace
		if ev.Keychar == 8 {
			return KeyBackspace
		}
		// For other printable characters, return the character
		return char
	default:
		// Check if keyStr is a special character that should be mapped
		runes := []rune(keyStr)
		if len(runes) == 1 {
			charFromKeyStr := string(runes[0])
			// Check for special keys by their character representation
			if charFromKeyStr == " " {
				log.Printf("[DEBUG] translateEvent: default case - mapping space character to KeySpace")
				return KeySpace
			}
			if charFromKeyStr == "\r" || charFromKeyStr == "\n" {
				log.Printf("[DEBUG] translateEvent: default case - mapping newline to KeyEnter")
				return KeyEnter
			}
			if charFromKeyStr == "\t" {
				log.Printf("[DEBUG] translateEvent: default case - mapping tab to KeyTab")
				return KeyTab
			}
			// For other printable characters, return the character
			return charFromKeyStr
		}
		// If keyStr is empty or multi-character, check keychar directly
		if keyStr == "" {
			// keyStr is empty, check keychar
			if char == " " {
				log.Printf("[DEBUG] translateEvent: empty keyStr, keychar is space, returning KeySpace")
				return KeySpace
			}
			if char == "\r" || char == "\n" {
				log.Printf("[DEBUG] translateEvent: empty keyStr, keychar is enter, returning KeyEnter")
				return KeyEnter
			}
			if char == "\t" {
				log.Printf("[DEBUG] translateEvent: empty keyStr, keychar is tab, returning KeyTab")
				return KeyTab
			}
		}
		// Ignore other special keys (arrows, function keys, etc.) here.
		return ""
	}
}
