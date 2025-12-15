package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/getlantern/systray"
	"github.com/go-vgo/robotgo"

	"text-expander/config"
	"text-expander/expander"
	"text-expander/gui"
	"text-expander/utils"
)

const version = "0.1.0"

// Global Fyne app instance for GUI windows
var fyneApp fyne.App

func main() {
	if err := os.MkdirAll("config", 0o755); err != nil {
		log.Fatalf("failed to create config directory: %v", err)
	}
	if err := os.MkdirAll("logs", 0o755); err != nil {
		log.Fatalf("failed to create logs directory: %v", err)
	}

	cfgPath := defaultConfigPath()
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	exp := expander.NewExpander(cfg)
	logger := utils.NewLogger(defaultLogPath())
	exp.SetLogger(logger)

	// Initialize Fyne app before systray
	fyneApp = app.NewWithID("com.textexpander.manager")

	systray.Run(func() { onReady(exp, cfg, logger) }, func() { onExit(exp, logger) })
}

func defaultConfigPath() string {
	return filepath.Join("config", "expansions.json")
}

func defaultLogPath() string {
	return filepath.Join("logs", "expander.log")
}

func onReady(exp *expander.Expander, cfg *config.Config, logger *utils.Logger) {
	// Load custom icon if available
	loadIcon()

	systray.SetTitle("Text Expander")
	systray.SetTooltip("Text Expander")

	// Set up notification callback
	exp.SetNotificationCallback(gui.ShowExpansionNotification)

	// Check for first run and show welcome dialog
	go checkFirstRun()

	if err := exp.Start(); err != nil {
		log.Printf("failed to start keyboard hook: %v", err)
	} else {
		log.Printf("keyboard hook started successfully")
	}

	toggleItem := systray.AddMenuItem("Disable", "Enable or disable expansions")
	configureItem := systray.AddMenuItem("Configure...", "Open configuration editor")
	statsItem := systray.AddMenuItem("Statistics", "Show expansion statistics")
	viewLogsItem := systray.AddMenuItem("View Logs", "Open log file")
	reloadItem := systray.AddMenuItem("Reload Configuration", "Reload expansions configuration")
	aboutItem := systray.AddMenuItem("About", "About Text Expander")

	systray.AddSeparator()
	quitItem := systray.AddMenuItem("Quit", "Quit the application")

	updateTrayTooltip(cfg)
	updateToggleTitle(cfg, toggleItem)

	go func() {
		for {
			select {
			case <-toggleItem.ClickedCh:
				toggleEnabled(cfg, toggleItem)
				updateTrayTooltip(cfg)
			case <-configureItem.ClickedCh:
				// Launch GUI as separate process
				go func() {
					defer func() {
						if r := recover(); r != nil {
							log.Printf("Error launching GUI: %v", r)
						}
					}()

					// Get executable directory
					exePath, err := os.Executable()
					if err != nil {
						log.Printf("Failed to get executable path: %v", err)
						robotgo.Alert("Error", "Failed to open configuration window")
						return
					}
					exeDir := filepath.Dir(exePath)
					guiPath := filepath.Join(exeDir, "gui-config.exe")

					// Launch the separate GUI executable
					cmd := exec.Command(guiPath)
					cmd.Dir = exeDir // Set working directory
					if err := cmd.Start(); err != nil {
						log.Printf("Failed to launch GUI: %v", err)
						robotgo.Alert("Error", fmt.Sprintf("Failed to open configuration window: %v", err))
					}
				}()
			case <-statsItem.ClickedCh:
				showStats(logger)
			case <-viewLogsItem.ClickedCh:
				openLogFile(defaultLogPath())
			case <-reloadItem.ClickedCh:
				exp.ReloadConfig()
			case <-aboutItem.ClickedCh:
				showAbout()
			case <-quitItem.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func onExit(exp *expander.Expander, logger *utils.Logger) {
	exp.Stop()
	if logger != nil {
		logger.Close()
	}
}

func toggleEnabled(cfg *config.Config, item *systray.MenuItem) {
	s := cfg.GetSettings()
	s.Enabled = !s.Enabled
	cfg.UpdateSettings(s)
	_ = cfg.Save()
	updateToggleTitle(cfg, item)
}

func updateToggleTitle(cfg *config.Config, item *systray.MenuItem) {
	s := cfg.GetSettings()
	if s.Enabled {
		item.SetTitle("Disable")
	} else {
		item.SetTitle("Enable")
	}
}

func updateTrayTooltip(cfg *config.Config) {
	s := cfg.GetSettings()
	if s.Enabled {
		systray.SetTooltip("Text Expander (Enabled)")
	} else {
		systray.SetTooltip("Text Expander (Disabled)")
	}
}

func showStats(logger *utils.Logger) {
	if logger == nil {
		return
	}
	stats := logger.GetStats()
	msg := fmt.Sprintf(
		"Total expansions: %d\nToday's expansions: %d\nMost used trigger: %s\nLast expansion: %s",
		stats.TotalExpansions,
		stats.TodayExpansions,
		stats.MostUsedTrigger,
		formatTime(stats.LastExpansion),
	)
	robotgo.Alert("Text Expander Statistics", msg)
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return "N/A"
	}
	return t.Format("2006-01-02 15:04:05")
}

func openLogFile(path string) {
	if _, err := os.Stat(path); err != nil {
		return
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", path)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", path)
	default:
		cmd = exec.Command("xdg-open", path)
	}

	_ = cmd.Start()
}

func showAbout() {
	msg := "Text Expander\n" +
		"Version " + version + "\n\n" +
		"Text Expander - A productivity tool"
	robotgo.Alert("About Text Expander", msg)
}

func loadIcon() {
	// Try to load custom icon
	iconPath := "app-icon.ico"
	if _, err := os.Stat(iconPath); err == nil {
		// Icon file exists, load it
		if iconData, err := os.ReadFile(iconPath); err == nil {
			systray.SetIcon(iconData)
			log.Printf("Loaded custom icon from %s", iconPath)
			return
		}
	}

	// Fallback: try PNG icon
	iconPath = "app-icon.png"
	if _, err := os.Stat(iconPath); err == nil {
		if iconData, err := os.ReadFile(iconPath); err == nil {
			systray.SetIcon(iconData)
			log.Printf("Loaded custom icon from %s (PNG)", iconPath)
			return
		}
	}

	log.Printf("No custom icon found, using default")
}

func checkFirstRun() {
	isFirst, err := utils.IsFirstRun()
	if err != nil {
		log.Printf("Error checking first run: %v", err)
		return
	}

	if isFirst {
		log.Printf("First run detected, showing welcome message")

		// Welcome message
		time.Sleep(2 * time.Second) // Wait for systray to initialize
		robotgo.Alert("Welcome to Text Expander!",
			"Text Expander is now running in your system tray (bottom-right corner).\\n\\n"+
				"Try typing: ;hello followed by Space\\n"+
				"You'll see it expand to 'Hello, World!'\\n\\n"+
				"Right-click the tray icon to configure expansions.\\n\\n"+
				"144+ built-in expansions ready to use!")

		utils.MarkAsCompleted()

		log.Printf("First-run setup completed")
	}
}
