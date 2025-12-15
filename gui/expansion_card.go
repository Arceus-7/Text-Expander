package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"text-expander/config"
)

// ExpansionCard represents a visual card for displaying an expansion
type ExpansionCard struct {
	widget.BaseWidget
	expansion *config.Expansion
	onEdit    func(*config.Expansion)
	onDelete  func(string)
}

// NewExpansionCard creates a new expansion card widget
func NewExpansionCard(exp *config.Expansion, onEdit func(*config.Expansion), onDelete func(string)) *ExpansionCard {
	card := &ExpansionCard{
		expansion: exp,
		onEdit:    onEdit,
		onDelete:  onDelete,
	}
	card.ExtendBaseWidget(card)
	return card
}

// CreateRenderer creates the visual representation of the card
func (c *ExpansionCard) CreateRenderer() fyne.WidgetRenderer {
	// Trigger text (large, bold)
	trigger := canvas.NewText(c.expansion.Trigger, theme.ForegroundColor())
	trigger.TextSize = 16
	trigger.TextStyle = fyne.TextStyle{Bold: true}

	// Description text (subtitle)
	description := canvas.NewText(c.expansion.Description, theme.DisabledTextColor())
	description.TextSize = 12

	// Preview text (truncated replacement)
	preview := c.getTruncatedPreview()
	previewText := canvas.NewText(preview, theme.ForegroundColor())
	previewText.TextSize = 11
	previewText.TextStyle = fyne.TextStyle{Monospace: true}

	// Category badge
	categoryBadge := c.createCategoryBadge()

	// Edit button
	editBtn := widget.NewButton("Edit", func() {
		if c.onEdit != nil {
			c.onEdit(c.expansion)
		}
	})
	editBtn.Importance = widget.LowImportance

	// Delete button
	deleteBtn := widget.NewButton("Delete", func() {
		if c.onDelete != nil {
			c.onDelete(c.expansion.Trigger)
		}
	})
	deleteBtn.Importance = widget.DangerImportance

	// Card background
	bg := canvas.NewRectangle(color.NRGBA{R: 245, G: 245, B: 245, A: 255})

	// Layout the card
	topRow := container.NewBorder(nil, nil, trigger, categoryBadge)
	middleRow := container.NewVBox(description, previewText)
	bottomRow := container.NewHBox(editBtn, deleteBtn)

	content := container.NewBorder(topRow, bottomRow, nil, nil, middleRow)

	cardWithBg := container.NewStack(bg, container.NewPadded(content))

	return &expansionCardRenderer{
		card:    c,
		bg:      bg,
		objects: []fyne.CanvasObject{cardWithBg},
	}
}

// getTruncatedPreview returns a truncated version of the replacement text
func (c *ExpansionCard) getTruncatedPreview() string {
	replacement := c.expansion.Replacement
	maxLen := 60

	// Replace newlines with spaces for preview
	preview := ""
	for _, char := range replacement {
		if char == '\n' || char == '\r' {
			preview += " "
		} else {
			preview += string(char)
		}
	}

	if len(preview) > maxLen {
		return preview[:maxLen] + "..."
	}
	return preview
}

// createCategoryBadge creates a colored badge for the category
func (c *ExpansionCard) createCategoryBadge() *fyne.Container {
	if c.expansion.Category == "" {
		return container.NewHBox()
	}

	categoryColor := GetCategoryColor(c.expansion.Category)

	badge := canvas.NewRectangle(categoryColor)
	badge.SetMinSize(fyne.NewSize(80, 20))

	label := canvas.NewText(c.expansion.Category, color.White)
	label.TextSize = 10
	label.Alignment = fyne.TextAlignCenter

	return container.NewStack(badge, container.NewCenter(label))
}

// expansionCardRenderer handles the rendering of the card
type expansionCardRenderer struct {
	card    *ExpansionCard
	bg      *canvas.Rectangle
	objects []fyne.CanvasObject
}

func (r *expansionCardRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
	for _, obj := range r.objects {
		obj.Resize(size)
	}
}

func (r *expansionCardRenderer) MinSize() fyne.Size {
	return fyne.NewSize(400, 100)
}

func (r *expansionCardRenderer) Refresh() {
	canvas.Refresh(r.bg)
	for _, obj := range r.objects {
		canvas.Refresh(obj)
	}
}

func (r *expansionCardRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *expansionCardRenderer) Destroy() {}
