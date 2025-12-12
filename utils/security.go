package utils

import (
	"strings"
	"sync"
	"time"

	"github.com/go-vgo/robotgo"
)

// Simple heuristics for determining when expansions should be disabled for
// security reasons, plus a rate limiter to avoid accidental rapid-fire
// expansions.

var (
	rateMu            sync.Mutex
	recentExpansions  []time.Time
	rateWindow        = 2 * time.Second
	maxExpansionsInWindow = 20
)

var blacklistedAppSubstrings = []string{
	"1password",
	"lastpass",
	"keepass",
	"bitwarden",
	"keychain",
	"authy",
	"dashlane",
}

// IsPasswordField attempts to determine whether the active window is a password
// prompt. This is a heuristic based on the window title and may not be perfect,
// but it errs on the side of disabling expansions.
func IsPasswordField() bool {
	title := strings.ToLower(robotgo.GetTitle())
	if strings.Contains(title, "password") || strings.Contains(title, "passcode") {
		return true
	}
	return false
}

// IsBlacklistedApp reports whether the current active window appears to belong
// to a blacklisted application (such as password managers).
func IsBlacklistedApp() bool {
	title := strings.ToLower(robotgo.GetTitle())
	for _, s := range blacklistedAppSubstrings {
		if strings.Contains(title, s) {
			return true
		}
	}
	return false
}

// ShouldAllowExpansion combines password detection, application blacklisting,
// and rate limiting to decide whether an expansion should proceed.
func ShouldAllowExpansion() bool {
	if IsPasswordField() || IsBlacklistedApp() {
		return false
	}

	rateMu.Lock()
	defer rateMu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rateWindow)

	// Drop old timestamps.
	filtered := recentExpansions[:0]
	for _, ts := range recentExpansions {
		if ts.After(cutoff) {
			filtered = append(filtered, ts)
		}
	}
	recentExpansions = filtered

	if len(recentExpansions) >= maxExpansionsInWindow {
		return false
	}

	recentExpansions = append(recentExpansions, now)
	return true
}