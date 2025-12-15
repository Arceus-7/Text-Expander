package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
)

// Expansion defines a single trigger/replacement pair as stored in the config file.
type Expansion struct {
	Trigger       string `json:"trigger"`
	Replacement   string `json:"replacement"`
	CaseSensitive bool   `json:"case_sensitive"`
	Description   string `json:"description"`
	Category      string `json:"category,omitempty"` // NEW: Category for filtering/organization
}

// Settings contains global behaviour flags.
type Settings struct {
	Enabled           bool `json:"enabled"`
	TriggerOnSpace    bool `json:"trigger_on_space"`
	TriggerOnTab      bool `json:"trigger_on_tab"`
	TriggerOnEnter    bool `json:"trigger_on_enter"`
	ShowNotifications bool `json:"show_notifications"`
	LogExpansions     bool `json:"log_expansions"`
}

// Config is the root configuration object for the application.
type Config struct {
	Expansions      []Expansion       `json:"expansions"`
	CustomVariables map[string]string `json:"custom_variables"`
	Settings        Settings          `json:"settings"`

	filePath string
	mu       sync.RWMutex
}

// LoadConfig loads configuration from the given path. If the file does not
// exist, it is created with a default configuration. If the file is corrupted,
// it is backed up and replaced with defaults.
func LoadConfig(path string) (*Config, error) {
	if path == "" {
		return nil, errors.New("config path is required")
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("creating config dir: %w", err)
	}

	cfg := &Config{}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Create default config.
			cfg = defaultConfig()
			cfg.filePath = path
			if err := cfg.Save(); err != nil {
				return cfg, fmt.Errorf("saving default config: %w", err)
			}
			return cfg, nil
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}

	if len(data) == 0 {
		cfg = defaultConfig()
	} else if err := json.Unmarshal(data, cfg); err != nil {
		// Backup corrupted file and fall back to defaults.
		backup := path + ".bak"
		_ = os.WriteFile(backup, data, 0o600)
		cfg = defaultConfig()
	}

	if cfg.CustomVariables == nil {
		cfg.CustomVariables = make(map[string]string)
	}

	cfg.filePath = path
	return cfg, nil
}

// Save writes the configuration to disk atomically.
func (c *Config) Save() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.filePath == "" {
		return errors.New("config file path is not set")
	}

	tmpPath := c.filePath + ".tmp"

	out := struct {
		Expansions      []Expansion       `json:"expansions"`
		CustomVariables map[string]string `json:"custom_variables"`
		Settings        Settings          `json:"settings"`
	}{
		Expansions:      c.Expansions,
		CustomVariables: c.CustomVariables,
		Settings:        c.Settings,
	}

	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	if err := os.WriteFile(tmpPath, data, 0o644); err != nil {
		return fmt.Errorf("write temp config: %w", err)
	}

	if err := os.Rename(tmpPath, c.filePath); err != nil {
		return fmt.Errorf("rename temp config: %w", err)
	}

	return nil
}

// AddExpansion adds a new expansion to the configuration, ensuring that the
// trigger is not empty and not duplicated.
func (c *Config) AddExpansion(exp Expansion) error {
	if strings.TrimSpace(exp.Trigger) == "" {
		return errors.New("trigger cannot be empty")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	for _, existing := range c.Expansions {
		if existing.Trigger == exp.Trigger {
			return fmt.Errorf("expansion with trigger %q already exists", exp.Trigger)
		}
	}

	c.Expansions = append(c.Expansions, exp)
	return nil
}

// RemoveExpansion removes an expansion by its trigger.
func (c *Config) RemoveExpansion(trigger string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	index := -1
	for i, exp := range c.Expansions {
		if exp.Trigger == trigger {
			index = i
			break
		}
	}
	if index == -1 {
		return fmt.Errorf("expansion with trigger %q not found", trigger)
	}

	c.Expansions = append(c.Expansions[:index], c.Expansions[index+1:]...)
	return nil
}

// Watch sets up a file watcher on the configuration file and calls the
// callback whenever the file is modified or recreated. The callback is called
// from a background goroutine.
func (c *Config) Watch(callback func()) error {
	c.mu.RLock()
	path := c.filePath
	c.mu.RUnlock()

	if path == "" {
		return errors.New("config file path is not set")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("create watcher: %w", err)
	}

	if err := watcher.Add(path); err != nil {
		_ = watcher.Close()
		return fmt.Errorf("watch config: %w", err)
	}

	go func() {
		defer watcher.Close()

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename) != 0 {
					if callback != nil {
						callback()
					}
				}
			case _, ok := <-watcher.Errors:
				if !ok {
					return
				}
				// Errors are ignored here; caller can provide separate logging if desired.
			}
		}
	}()

	return nil
}

// ConfigPath returns the underlying configuration file path.
func (c *Config) ConfigPath() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.filePath
}

// GetExpansions returns a copy of the expansion slice.
func (c *Config) GetExpansions() []Expansion {
	c.mu.RLock()
	defer c.mu.RUnlock()

	exps := make([]Expansion, len(c.Expansions))
	copy(exps, c.Expansions)
	return exps
}

// GetSettings returns a copy of the current settings.
func (c *Config) GetSettings() Settings {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Settings
}

// UpdateSettings replaces the current settings with the provided value.
func (c *Config) UpdateSettings(s Settings) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Settings = s
}

// GetCustomVars returns a copy of the custom variable map.
func (c *Config) GetCustomVars() map[string]string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cp := make(map[string]string, len(c.CustomVariables))
	for k, v := range c.CustomVariables {
		cp[k] = v
	}
	return cp
}

// SetCustomVar sets a custom variable value.
func (c *Config) SetCustomVar(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.CustomVariables == nil {
		c.CustomVariables = make(map[string]string)
	}
	c.CustomVariables[key] = value
}

// DeleteCustomVar removes a custom variable.
func (c *Config) DeleteCustomVar(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.CustomVariables, key)
}

// defaultConfig returns a Config populated with sensible defaults.
func defaultConfig() *Config {
	return &Config{
		Expansions: []Expansion{
			{
				Trigger:       ";email",
				Replacement:   "your.email@example.com",
				CaseSensitive: false,
				Description:   "Personal email",
			},
			{
				Trigger:       ";date",
				Replacement:   "{DATE}",
				CaseSensitive: false,
				Description:   "Current date",
			},
			{
				Trigger:       ";sig",
				Replacement:   "Best regards,\nYour Name",
				CaseSensitive: false,
				Description:   "Email signature",
			},
			{
				Trigger:       ";shrug",
				Replacement:   "¯\\_(ツ)_/¯",
				CaseSensitive: false,
				Description:   "Shrug emoji",
			},
		},
		CustomVariables: map[string]string{
			"NAME":    "John Doe",
			"COMPANY": "Acme Corp",
		},
		Settings: Settings{
			Enabled:           true,
			TriggerOnSpace:    true,
			TriggerOnTab:      true,
			TriggerOnEnter:    true,
			ShowNotifications: false,
			LogExpansions:     true,
		},
	}
}
