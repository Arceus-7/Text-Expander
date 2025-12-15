package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type AppSettings struct {
	FirstRun          bool      `json:"first_run"`
	LastVersion       string    `json:"last_version"`
	ShowNotifications bool      `json:"show_notifications"`
	NotificationDur   int       `json:"notification_duration"` // seconds
	PlaySound         bool      `json:"play_sound"`
	UserName          string    `json:"user_name"`
	UserEmail         string    `json:"user_email"`
	InstallDate       time.Time `json:"install_date"`
	LastOpenDate      time.Time `json:"last_open_date"`
}

const settingsFile = "config/app_settings.json"
const currentVersion = "1.1.0"

// LoadSettings loads app settings from disk
func LoadSettings() (*AppSettings, error) {
	// If settings file doesn't exist, this is first run
	if _, err := os.Stat(settingsFile); os.IsNotExist(err) {
		return &AppSettings{
			FirstRun:          true,
			LastVersion:       currentVersion,
			ShowNotifications: true,
			NotificationDur:   2,
			PlaySound:         false,
			InstallDate:       time.Now(),
			LastOpenDate:      time.Now(),
		}, nil
	}

	data, err := os.ReadFile(settingsFile)
	if err != nil {
		return nil, err
	}

	var settings AppSettings
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, err
	}

	settings.LastOpenDate = time.Now()
	return &settings, nil
}

// SaveSettings saves app settings to disk
func SaveSettings(settings *AppSettings) error {
	// Ensure config directory exists
	if err := os.MkdirAll(filepath.Dir(settingsFile), 0755); err != nil {
		return err
	}

	settings.LastVersion = currentVersion
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(settingsFile, data, 0644)
}

// IsFirstRun returns true if this is the first time the app is run
func IsFirstRun() (bool, error) {
	settings, err := LoadSettings()
	if err != nil {
		return false, err
	}
	return settings.FirstRun, nil
}

// MarkAsCompleted marks the first run as completed
func MarkAsCompleted() error {
	settings, err := LoadSettings()
	if err != nil {
		return err
	}

	settings.FirstRun = false
	return SaveSettings(settings)
}

// UpdateUserInfo updates user personalization info
func UpdateUserInfo(name, email string) error {
	settings, err := LoadSettings()
	if err != nil {
		return err
	}

	settings.UserName = name
	settings.UserEmail = email
	return SaveSettings(settings)
}

// GetSettings returns the current settings
func GetSettings() (*AppSettings, error) {
	return LoadSettings()
}
