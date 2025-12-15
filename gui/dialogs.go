package gui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"text-expander/config"
)

// ShowExpansionDialog shows a simplified, readable dialog
func ShowExpansionDialog(parent fyne.Window, cfg *config.Config, existing *config.Expansion, onSave func()) {
	// Colors
	labelColor := color.NRGBA{R: 99, G: 102, B: 241, A: 255}
	hintColor := color.NRGBA{R: 107, G: 114, B: 128, A: 255}

	// Fields
	triggerEntry := widget.NewEntry()
	triggerEntry.SetPlaceHolder("e.g., ;myshortcut")

	descEntry := widget.NewEntry()
	descEntry.SetPlaceHolder("e.g., My custom shortcut")

	replacementEntry := widget.NewMultiLineEntry()
	replacementEntry.SetPlaceHolder("What the trigger expands to...")
	replacementEntry.SetMinRowsVisible(4)

	categories := []string{"Personal", "Python", "JavaScript", "Go", "C", "SQL", "HTML", "CSS", "Professional", "Symbols", "General"}
	categorySelect := widget.NewSelect(categories, nil)
	categorySelect.PlaceHolder = "Select category..."

	caseSensitiveCheck := widget.NewCheck("Case sensitive", nil)

	// Fill existing
	isEdit := existing != nil
	if isEdit {
		triggerEntry.SetText(existing.Trigger)
		descEntry.SetText(existing.Description)
		replacementEntry.SetText(existing.Replacement)
		categorySelect.SetSelected(existing.Category)
		caseSensitiveCheck.SetChecked(existing.CaseSensitive)
	}

	// Labels
	triggerLabel := canvas.NewText("Trigger (what you type)", labelColor)
	triggerLabel.TextSize = 16
	triggerLabel.TextStyle = fyne.TextStyle{Bold: true}

	triggerHint := canvas.NewText("Start with ; or your preferred prefix", hintColor)
	triggerHint.TextSize = 12

	descLabel := canvas.NewText("Description (what it does)", labelColor)
	descLabel.TextSize = 16
	descLabel.TextStyle = fyne.TextStyle{Bold: true}

	replLabel := canvas.NewText("Replacement (what it becomes)", labelColor)
	replLabel.TextSize = 16
	replLabel.TextStyle = fyne.TextStyle{Bold: true}

	replHint := canvas.NewText("Use {DATE}, {TIME}, {CURSOR}, {CLIPBOARD} as variables", hintColor)
	replHint.TextSize = 12

	categoryLabel := canvas.NewText("Category", labelColor)
	categoryLabel.TextSize = 16
	categoryLabel.TextStyle = fyne.TextStyle{Bold: true}

	optionsLabel := canvas.NewText("Options", labelColor)
	optionsLabel.TextSize = 16
	optionsLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Form
	form := container.NewVBox(
		triggerLabel,
		triggerEntry,
		triggerHint,
		widget.NewLabel(""),

		descLabel,
		descEntry,
		widget.NewLabel(""),

		replLabel,
		replacementEntry,
		replHint,
		widget.NewLabel(""),

		categoryLabel,
		categorySelect,
		widget.NewLabel(""),

		optionsLabel,
		caseSensitiveCheck,
	)

	// White background
	bg := canvas.NewRectangle(color.NRGBA{R: 255, G: 255, B: 255, A: 255})
	formWithBg := container.NewStack(bg, container.NewPadded(form))

	// Scrollable
	scroll := container.NewVScroll(formWithBg)
	scroll.SetMinSize(fyne.NewSize(550, 550))

	// Validate
	validate := func() bool {
		if triggerEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("trigger cannot be empty"), parent)
			return false
		}
		if replacementEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("replacement cannot be empty"), parent)
			return false
		}
		return true
	}

	// Save
	saveFunc := func() {
		if !validate() {
			return
		}

		expansion := config.Expansion{
			Trigger:       triggerEntry.Text,
			Replacement:   replacementEntry.Text,
			Description:   descEntry.Text,
			Category:      categorySelect.Selected,
			CaseSensitive: caseSensitiveCheck.Checked,
		}

		if isEdit {
			cfg.RemoveExpansion(existing.Trigger)
		}

		if err := cfg.AddExpansion(expansion); err != nil {
			dialog.ShowError(err, parent)
			return
		}

		if err := cfg.Save(); err != nil {
			dialog.ShowError(err, parent)
			return
		}

		if onSave != nil {
			onSave()
		}
	}

	// Dialog
	title := "Add New Expansion"
	if isEdit {
		title = "Edit Expansion"
	}

	d := dialog.NewCustomConfirm(title, "Save", "Cancel", scroll, func(save bool) {
		if save {
			saveFunc()
		}
	}, parent)

	d.Resize(fyne.NewSize(600, 650))
	d.Show()
}

// ShowDeleteConfirmation shows a confirmation dialog before deleting
func ShowDeleteConfirmation(parent fyne.Window, trigger string, onConfirm func()) {
	dialog.ShowConfirm(
		"Delete Expansion",
		"Are you sure you want to delete '"+trigger+"'?\n\nThis action cannot be undone.",
		func(confirmed bool) {
			if confirmed && onConfirm != nil {
				onConfirm()
			}
		},
		parent,
	)
}

// ShowHelpDialog shows help with proper colors
func ShowHelpDialog(parent fyne.Window) {
	bg := canvas.NewRectangle(color.NRGBA{R: 255, G: 255, B: 255, A: 255})

	textColor := color.NRGBA{R: 31, G: 41, B: 55, A: 255}

	titleText := canvas.NewText("Text Expander Help", textColor)
	titleText.TextSize = 24
	titleText.TextStyle = fyne.TextStyle{Bold: true}

	howToTitle := canvas.NewText("HOW TO USE:", color.NRGBA{R: 99, G: 102, B: 241, A: 255})
	howToTitle.TextSize = 16
	howToTitle.TextStyle = fyne.TextStyle{Bold: true}

	howTo1 := canvas.NewText("1. Add expansions using the \"+ New Expansion\" button", textColor)
	howTo1.TextSize = 14
	howTo2 := canvas.NewText("2. Type a trigger (e.g., ;hello) and press Space/Tab/Enter", textColor)
	howTo2.TextSize = 14
	howTo3 := canvas.NewText("3. Watch it expand into your text!", textColor)
	howTo3.TextSize = 14

	varsTitle := canvas.NewText("TEMPLATE VARIABLES:", color.NRGBA{R: 99, G: 102, B: 241, A: 255})
	varsTitle.TextSize = 16
	varsTitle.TextStyle = fyne.TextStyle{Bold: true}

	var1 := canvas.NewText("• {DATE} - Current date (2024-12-15)", textColor)
	var1.TextSize = 14
	var2 := canvas.NewText("• {TIME} - Current time (21:30:45)", textColor)
	var2.TextSize = 14
	var3 := canvas.NewText("• {DATETIME} - Date and time combined", textColor)
	var3.TextSize = 14
	var4 := canvas.NewText("• {CLIPBOARD} - Paste clipboard content", textColor)
	var4.TextSize = 14
	var5 := canvas.NewText("• {CURSOR} - Position cursor after expansion", textColor)
	var5.TextSize = 14

	tipsTitle := canvas.NewText("TIPS:", color.NRGBA{R: 99, G: 102, B: 241, A: 255})
	tipsTitle.TextSize = 16
	tipsTitle.TextStyle = fyne.TextStyle{Bold: true}

	tip1 := canvas.NewText("✓ Use a consistent prefix like ; for triggers", textColor)
	tip1.TextSize = 14
	tip2 := canvas.NewText("✓ Keep triggers short and memorable", textColor)
	tip2.TextSize = 14
	tip3 := canvas.NewText("✓ Test new expansions in Notepad first", textColor)
	tip3.TextSize = 14
	tip4 := canvas.NewText("✓ Back up your config regularly", textColor)
	tip4.TextSize = 14

	spacer := canvas.NewText("", textColor)
	spacer.TextSize = 8

	content := container.NewVBox(
		titleText,
		spacer,
		howToTitle, howTo1, howTo2, howTo3,
		spacer,
		varsTitle, var1, var2, var3, var4, var5,
		spacer,
		tipsTitle, tip1, tip2, tip3, tip4,
	)

	paddedContent := container.NewPadded(content)
	contentWithBg := container.NewStack(bg, paddedContent)

	scroll := container.NewVScroll(contentWithBg)
	scroll.SetMinSize(fyne.NewSize(600, 500))

	d := dialog.NewCustom("Help", "Close", scroll, parent)
	d.Resize(fyne.NewSize(650, 600))
	d.Show()
}
