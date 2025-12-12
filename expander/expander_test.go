package expander

import (
	"testing"

	"github.com/yourusername/text-expander/config"
)

func TestMatchExpansionPrefersLongest(t *testing.T) {
	exps := map[string]Expansion{
		";e": {
			Trigger:      ";e",
			Replacement:  "short",
			CaseSensitive: false,
		},
		";email": {
			Trigger:      ";email",
			Replacement:  "long",
			CaseSensitive: false,
		},
	}

	exp, ok := matchExpansion("test;email", exps)
	if !ok {
		t.Fatalf("expected a matching expansion")
	}
	if exp.Trigger != ";email" {
		t.Fatalf("expected longest trigger ';email', got %q", exp.Trigger)
	}
}

func TestMatchExpansionCaseInsensitive(t *testing.T) {
	exps := map[string]Expansion{
		";date": {
			Trigger:      ";date",
			Replacement:  "{DATE}",
			CaseSensitive: false,
		},
	}

	exp, ok := matchExpansion("Today is ;DATE", exps)
	if !ok {
		t.Fatalf("expected a matching expansion for case-insensitive trigger")
	}
	if exp.Trigger != ";date" {
		t.Fatalf("unexpected trigger %q", exp.Trigger)
	}
}

func TestNewExpanderInitialisesFromConfig(t *testing.T) {
	cfg := &config.Config{
		Expansions: []config.Expansion{
			{Trigger: ";x", Replacement: "X"},
		},
		CustomVariables: map[string]string{
			"NAME": "Tester",
		},
	}

	e := NewExpanderWithKeyboard(cfg, &KeyboardHook{})
	if e == nil {
		t.Fatalf("expected non-nil expander")
	}

	if _, ok := e.expansions[";x"]; !ok {
		t.Fatalf("expected expansion ';x' to be loaded into expander")
	}
}