package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"text-expander/config"
)

// ShowExpansionDialog shows an enhanced dialog for adding or editing an expansion
func ShowExpansionDialog(parent fyne.Window, cfg *config.Config, existing *config.Expansion, onSave func()) {
	// Create form fields
	triggerEntry := widget.NewEntry()
	triggerEntry.SetPlaceHolder("e.g., ;myshortcut")

	descEntry := widget.NewEntry()
	descEntry.SetPlaceHolder("e.g., My custom shortcut")

	replacementEntry := widget.NewMultiLineEntry()
	replacementEntry.SetPlaceHolder("What the trigger expands to...")
	replacementEntry.SetMinRowsVisible(5)

	// Category selection
	categories := []string{"Personal", "Python", "JavaScript", "Go", "C", "SQL", "HTML", "CSS", "Professional", "Symbols", "General"}
	categorySelect := widget.NewSelect(categories, nil)
	categorySelect.PlaceHolder = "Select category..."

	caseSensitiveCheck := widget.NewCheck("Case sensitive", nil)

	// Fill existing data if editing
	isEdit := existing != nil
	if isEdit {
		triggerEntry.SetText(existing.Trigger)
		descEntry.SetText(existing.Description)
		replacementEntry.SetText(existing.Replacement)
		categorySelect.SetSelected(existing.Category)
		caseSensitiveCheck.SetChecked(existing.CaseSensitive)
	}

	// Template variable helper
	templateHelp := widget.NewLabel("Template Variables: {DATE} {TIME} {DATETIME} {CLIPBOARD} {CURSOR}")
	templateHelp.Wrapping = fyne.TextWrapWord

	// Add template buttons
	insertDate := widget.NewButton("{DATE}", func() {
		replacementEntry.SetText(replacementEntry.Text + "{DATE}")
	})
	insertTime := widget.NewButton("{TIME}", func() {
		replacementEntry.SetText(replacementEntry.Text + "{TIME}")
	})
	insertCursor := widget.NewButton("{CURSOR}", func() {
		replacementEntry.SetText(replacementEntry.Text + "{CURSOR}")
	})
	insertClipboard := widget.NewButton("{CLIPBOARD}", func() {
		replacementEntry.SetText(replacementEntry.Text + "{CLIPBOARD}")
	})

	templateButtons := container.NewHBox(insertDate, insertTime, insertCursor, insertClipboard)

	// Create form
	form := container.NewVBox(
		widget.NewLabelWithStyle("Trigger (what you type)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		triggerEntry,
		widget.NewLabel("Start with ; or your preferred prefix"),

		widget.NewSeparator(),

		widget.NewLabelWithStyle("Description (what it does)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		descEntry,

		widget.NewSeparator(),

		widget.NewLabelWithStyle("Replacement (what it becomes)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		replacementEntry,

		widget.NewLabel("Template Variables:"),
		templateButtons,
		templateHelp,

		widget.NewSeparator(),

		widget.NewLabelWithStyle("Category", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		categorySelect,

		widget.NewLabel("Options:"),
		caseSensitiveCheck,
	)

	// Validation function
	validate := func() bool {
		if triggerEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("Trigger cannot be empty"), parent)
			return false
		}
		if replacementEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("Replacement cannot be empty"), parent)
			return false
		}
		return true
	}

	// Save button handler
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
			// Remove old expansion
			cfg.RemoveExpansion(existing.Trigger)
		}

		// Add new/updated expansion
		if err := cfg.AddExpansion(expansion); err != nil {
			dialog.ShowError(err, parent)
			return
		}

		// Save to file
		if err := cfg.Save(); err != nil {
			dialog.ShowError(err, parent)
			return
		}

		if onSave != nil {
			onSave()
		}
	}

	// Create dialog
	title := "Add New Expansion"
	if isEdit {
		title = "Edit Expansion"
	}

	d := dialog.NewCustomConfirm(title, "Save", "Cancel", form, func(save bool) {
		if save {
			saveFunc()
		}
	}, parent)

	d.Resize(fyne.NewSize(600, 700))
	d.Show()
}

// ShowDeleteConfirmation shows a confirmation dialog before deleting an expansion
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

// ShowHelpDialog shows a help dialog with usage instructions
func ShowHelpDialog(parent fyne.Window) {
	helpText := `
Text Expander Help

HOW TO USE:
1. Add expansions using the "+ New" button
2. Type a trigger (e.g., ;hello) and press Space/Tab/Enter
3. Watch it expand into your text!

TEMPLATE VARIABLES:
- {DATE} - Current date (2024-12-15)
- {TIME} - Current time (21:30:45)
- {DATETIME} - Date and time combined
- {CLIPBOARD} - Paste clipboard content
- {CURSOR} - Position cursor after expansion

TIPS:
- Use a consistent prefix like ; for triggers
- Keep triggers short and memorable
- Test new expansions in Notepad first
- Use categories to organize expansions
- Back up your config regularly

CATEGORIES:
Organize expansions by type (Python, SQL, Personal, etc.) for easier management.

KEYBOARD SHORTCUTS:
- Ctrl+N - New expansion
- Ctrl+F - Focus search
- Delete - Delete selected expansion
- Esc - Close dialog
`

	d := dialog.NewInformation("Help", helpText, parent)
	d.Resize(fyne.NewSize(500, 600))
	d.Show()
}
