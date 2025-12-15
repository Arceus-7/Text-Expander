package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"text-expander/config"
)

// ExpansionCard is a beautiful card widget for displaying expansions
type ExpansionCard struct {
	widget.BaseWidget
	expansion  *config.Expansion
	onEdit     func(*config.Expansion)
	onDelete   func(string)
	background *canvas.Rectangle
	shadow     *canvas.Rectangle
	hovered    bool
}

// NewExpansionCard creates a new modern expansion card with shadows and hover effects
func NewExpansionCard(exp *config.Expansion, onEdit func(*config.Expansion), onDelete func(string)) *ExpansionCard {
	card := &ExpansionCard{
		expansion: exp,
		onEdit:    onEdit,
		onDelete:  onDelete,
	}
	card.ExtendBaseWidget(card)
	return card
}

// CreateRenderer creates the card renderer with modern styling
func (c *ExpansionCard) CreateRenderer() fyne.WidgetRenderer {
	// Shadow layer (underneath card)
	c.shadow = canvas.NewRectangle(ColorShadow)
	c.shadow.CornerRadius = 10

	// Main card background
	c.background = canvas.NewRectangle(ColorSurface)
	c.background.CornerRadius = 8

	// Trigger text (large, bold)
	triggerText := canvas.NewText(c.expansion.Trigger, ColorTextPrimary)
	triggerText.TextSize = 20 // Larger for emphasis
	triggerText.TextStyle = fyne.TextStyle{Bold: true}

	// Description text (smaller, lighter)
	descText := canvas.NewText(c.expansion.Description, ColorTextSecondary)
	descText.TextSize = 15

	// Preview text (smaller, monospace feel)
	preview := c.expansion.Replacement
	// Clean newlines for preview
	previewClean := ""
	for _, char := range preview {
		if char == '\n' || char == '\r' {
			previewClean += " "
		} else {
			previewClean += string(char)
		}
	}
	if len(previewClean) > 100 {
		previewClean = previewClean[:100] + "..."
	}
	previewText := canvas.NewText(previewClean, ColorTextTertiary)
	previewText.TextSize = 13
	previewText.TextStyle = fyne.TextStyle{Monospace: true}

	// Category badge (if exists)
	var categoryBadge fyne.CanvasObject
	if c.expansion.Category != "" {
		badgeColor := getCategoryColor(c.expansion.Category)
		badgeBg := canvas.NewRectangle(badgeColor)
		badgeBg.CornerRadius = 4
		badgeBg.Resize(fyne.NewSize(float32(len(c.expansion.Category)*8+16), 24))

		badgeLabel := canvas.NewText(c.expansion.Category, ColorSurface)
		badgeLabel.TextSize = 12
		badgeLabel.TextStyle = fyne.TextStyle{Bold: true}

		categoryBadge = container.NewStack(badgeBg, container.NewCenter(badgeLabel))
	} else {
		categoryBadge = layout.NewSpacer()
	}

	// Edit button (ghost style with emoji)
	editBtn := widget.NewButton("Edit", func() {
		if c.onEdit != nil {
			c.onEdit(c.expansion)
		}
	})
	editBtn.Importance = widget.LowImportance

	// Delete button (danger style with emoji)
	deleteBtn := widget.NewButton("Delete", func() {
		if c.onDelete != nil {
			c.onDelete(c.expansion.Trigger)
		}
	})
	deleteBtn.Importance = widget.DangerImportance

	// Button container
	buttons := container.NewHBox(
		layout.NewSpacer(),
		editBtn,
		deleteBtn,
	)

	// Main content layout with better spacing
	spacer1 := widget.NewLabel("")
	spacer2 := widget.NewLabel("")
	spacer3 := widget.NewLabel("")

	content := container.NewVBox(
		triggerText,
		spacer1,
		descText,
		spacer1,
		previewText,
		spacer2,
		categoryBadge,
		spacer3,
		buttons,
	)

	// Add generous padding
	paddedContent := container.NewPadded(container.NewPadded(content))

	// Stack: shadow -> background -> content
	objects := []fyne.CanvasObject{
		c.shadow,
		c.background,
		paddedContent,
	}

	return &expansionCardRenderer{
		card:    c,
		objects: objects,
	}
}

// Tapped handles tap events
func (c *ExpansionCard) Tapped(_ *fyne.PointEvent) {
	// Could add tap animation here
}

// TappedSecondary handles right-click
func (c *ExpansionCard) TappedSecondary(_ *fyne.PointEvent) {
}

// MouseIn handles mouse enter (hover effect)
func (c *ExpansionCard) MouseIn(_ *desktop.MouseEvent) {
	c.hovered = true
	// Strengthen shadow on hover
	c.shadow.FillColor = ColorShadowStrong
	c.shadow.Refresh()
}

// MouseOut handles mouse leave
func (c *ExpansionCard) MouseOut() {
	c.hovered = false
	// Return to normal shadow
	c.shadow.FillColor = ColorShadow
	c.shadow.Refresh()
}

// MouseMoved handles mouse movement
func (c *ExpansionCard) MouseMoved(_ *desktop.MouseEvent) {
}

// getCategoryColor returns color for category badges
func getCategoryColor(category string) color.Color {
	colors := map[string]color.Color{
		"Python":       color.NRGBA{R: 59, G: 130, B: 246, A: 255}, // Blue
		"JavaScript":   color.NRGBA{R: 251, G: 191, B: 36, A: 255}, // Yellow
		"Go":           color.NRGBA{R: 6, G: 182, B: 212, A: 255},  // Cyan
		"SQL":          color.NRGBA{R: 14, G: 165, B: 233, A: 255}, // Azure
		"C":            color.NRGBA{R: 75, G: 85, B: 99, A: 255},   // Gray
		"HTML":         color.NRGBA{R: 249, G: 115, B: 22, A: 255}, // Orange
		"CSS":          color.NRGBA{R: 37, G: 99, B: 235, A: 255},  // Royal Blue
		"Personal":     color.NRGBA{R: 168, G: 85, B: 247, A: 255}, // Purple
		"Professional": color.NRGBA{R: 34, G: 197, B: 94, A: 255},  // Green
		"Symbols":      color.NRGBA{R: 251, G: 146, B: 60, A: 255}, // Deep Orange
	}

	if c, ok := colors[category]; ok {
		return c
	}
	return ColorPrimary // Default
}

// expansionCardRenderer renders the card
type expansionCardRenderer struct {
	card    *ExpansionCard
	objects []fyne.CanvasObject
}

func (r *expansionCardRenderer) Destroy() {}

func (r *expansionCardRenderer) Layout(size fyne.Size) {
	// Position shadow slightly offset and larger
	shadowOffset := float32(2)
	r.card.shadow.Move(fyne.NewPos(shadowOffset, shadowOffset))
	r.card.shadow.Resize(size)

	// Position background
	r.card.background.Move(fyne.NewPos(0, 0))
	r.card.background.Resize(size)

	// Content fills the card
	r.objects[2].Resize(size)
}

func (r *expansionCardRenderer) MinSize() fyne.Size {
	return fyne.NewSize(450, 180) // Larger for better spacing
}

func (r *expansionCardRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *expansionCardRenderer) Refresh() {
	canvas.Refresh(r.card)
}
