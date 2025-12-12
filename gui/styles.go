package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// ApplyTheme applies a simple light theme for a clean, modern look.
func ApplyTheme(a fyne.App) {
	a.Settings().SetTheme(theme.LightTheme())
}