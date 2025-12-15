package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// Modern color palette
var (
	// Primary colors
	ColorPrimary     = color.NRGBA{R: 99, G: 102, B: 241, A: 255} // Indigo-500
	ColorPrimaryDark = color.NRGBA{R: 79, G: 70, B: 229, A: 255}  // Indigo-600
	ColorSecondary   = color.NRGBA{R: 139, G: 92, B: 246, A: 255} // Purple-500

	// Backgrounds
	ColorBackground = color.NRGBA{R: 250, G: 250, B: 250, A: 255} // Almost white
	ColorSurface    = color.NRGBA{R: 255, G: 255, B: 255, A: 255} // Pure white
	ColorHover      = color.NRGBA{R: 249, G: 250, B: 251, A: 255} // Gray-50

	// Text colors
	ColorTextPrimary   = color.NRGBA{R: 31, G: 41, B: 55, A: 255}    // Gray-800
	ColorTextSecondary = color.NRGBA{R: 107, G: 114, B: 128, A: 255} // Gray-500
	ColorTextTertiary  = color.NRGBA{R: 156, G: 163, B: 175, A: 255} // Gray-400

	// Borders
	ColorBorder      = color.NRGBA{R: 229, G: 231, B: 235, A: 255} // Gray-200
	ColorBorderLight = color.NRGBA{R: 243, G: 244, B: 246, A: 255} // Gray-100

	// Status colors
	ColorSuccess = color.NRGBA{R: 16, G: 185, B: 129, A: 255} // Green-500
	ColorWarning = color.NRGBA{R: 245, G: 158, B: 11, A: 255} // Amber-500
	ColorDanger  = color.NRGBA{R: 239, G: 68, B: 68, A: 255}  // Red-500

	// Shadows (used for layering)
	ColorShadowLight  = color.NRGBA{R: 0, G: 0, B: 0, A: 10}
	ColorShadow       = color.NRGBA{R: 0, G: 0, B: 0, A: 20}
	ColorShadowStrong = color.NRGBA{R: 0, G: 0, B: 0, A: 40}
)

// CustomTheme implements a modern, beautiful theme
type CustomTheme struct{}

var _ fyne.Theme = (*CustomTheme)(nil)

func (m *CustomTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return ColorBackground
	case theme.ColorNameButton:
		return ColorSurface
	case theme.ColorNameDisabled:
		return ColorTextTertiary
	case theme.ColorNameError:
		return ColorDanger
	case theme.ColorNameFocus:
		return ColorPrimary
	case theme.ColorNameForeground:
		return ColorTextPrimary
	case theme.ColorNameHover:
		return ColorHover
	case theme.ColorNameInputBackground:
		return ColorSurface
	case theme.ColorNameInputBorder:
		return ColorBorder
	case theme.ColorNameMenuBackground:
		return ColorSurface
	case theme.ColorNameOverlayBackground:
		return color.NRGBA{R: 0, G: 0, B: 0, A: 128}
	case theme.ColorNamePlaceHolder:
		return ColorTextSecondary
	case theme.ColorNamePressed:
		return ColorPrimaryDark
	case theme.ColorNamePrimary:
		return ColorPrimary
	case theme.ColorNameScrollBar:
		return ColorBorderLight
	case theme.ColorNameShadow:
		return ColorShadow
	case theme.ColorNameSuccess:
		return ColorSuccess
	case theme.ColorNameWarning:
		return ColorWarning
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

func (m *CustomTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m *CustomTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m *CustomTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return 16
	case theme.SizeNameInlineIcon:
		return 20
	case theme.SizeNameScrollBar:
		return 12
	case theme.SizeNameScrollBarSmall:
		return 8
	case theme.SizeNameSeparatorThickness:
		return 1
	case theme.SizeNameInputBorder:
		return 2
	case theme.SizeNameInputRadius:
		return 8
	case theme.SizeNameSelectionRadius:
		return 8
	default:
		return theme.DefaultTheme().Size(name)
	}
}

// ApplyTheme applies the custom modern theme to the app
func ApplyTheme(app fyne.App) {
	app.Settings().SetTheme(&CustomTheme{})
}
