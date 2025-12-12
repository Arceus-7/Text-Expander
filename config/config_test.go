package config

import (
	"path/filepath"
	"testing"
)

func TestLoadConfigCreatesDefault(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "expansions.json")

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}

	if len(cfg.Expansions) == 0 {
		t.Fatalf("expected default expansions to be created")
	}

	if cfg.ConfigPath() != path {
		t.Fatalf("expected config path %q, got %q", path, cfg.ConfigPath())
	}
}

func TestAddAndRemoveExpansion(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "expansions.json")

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}

	exp := Expansion{
		Trigger:     ";test",
		Replacement: "TEST",
	}

	if err := cfg.AddExpansion(exp); err != nil {
		t.Fatalf("AddExpansion returned error: %v", err)
	}

	found := false
	for _, e := range cfg.GetExpansions() {
		if e.Trigger == ";test" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected to find added expansion")
	}

	if err := cfg.RemoveExpansion(";test"); err != nil {
		t.Fatalf("RemoveExpansion returned error: %v", err)
	}
}