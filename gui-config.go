package main

import (
	"log"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"text-expander/config"
	"text-expander/gui"
)

func main() {
	// Get config path (same as main app)
	cfgPath := filepath.Join("config", "expansions.json")

	// Load configuration
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		return
	}

	// Create Fyne app
	a := app.NewWithID("com.textexpander.config")
	gui.ApplyTheme(a)

	// Create and show editor window
	w := a.NewWindow("Text Expander Manager")
	w.Resize(fyne.NewSize(1000, 700))
	w.CenterOnScreen()

	// Create editor state
	gui.CreateEditorWindow(w, cfg)

	// Show and run (blocks until window closed)
	w.ShowAndRun()
}
