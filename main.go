package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/getlantern/systray"
	"github.com/go-vgo/robotgo"

	"github.com/yourusername/text-expander/config"
	"github.com/yourusername/text-expander/expander"
	"github.com/yourusername/text-expander/gui"
	"github.com/yourusername/text-expander/utils"
)

const version = "0.1.0"

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

	systray.Run(func() { onReady(exp, cfg, logger) }, func() { onExit(exp, logger) })
}

func defaultConfigPath() string {
	return filepath.Join("config", "expansions.json")
}

func defaultLogPath() string {
	return filepath.Join("logs", "expander.log")
}

func onReady(exp *expander.Expander, cfg *config.Config, logger *utils.Logger) {
	systray.SetTitle("Text Expander")
	systray.SetTooltip("Text Expander")

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
				// Run editor in a separate goroutine - it creates its own Fyne app
				go func() {
					gui.ShowEditor(cfg)
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
		"https://github.com/yourusername/text-expander"
	robotgo.Alert("About Text Expander", msg)
}
