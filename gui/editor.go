package gui

import (
	"encoding/json"
	"fmt"
	"io"
	"runtime"
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"text-expander/config"
)

type editorState struct {
	cfg *config.Config

	window              fyne.Window
	searchEntry         *widget.Entry
	expansionsContainer *fyne.Container
	customVarsContainer *fyne.Container
	settingsContainer   *fyne.Container
	filteredExpansions  []config.Expansion
}

// ShowEditor launches the configuration editor as a standalone Fyne
// application. It operates directly on the provided Config instance.
// This function should be called from a goroutine as it blocks until
// the window is closed.
func ShowEditor(cfg *config.Config) {
	// Lock this goroutine to an OS thread to ensure Fyne operations
	// happen on the correct thread
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	a := app.New()
	ApplyTheme(a)

	w := a.NewWindow("Text Expander Configuration")
	w.Resize(fyne.NewSize(900, 600))

	state := &editorState{
		cfg:                 cfg,
		window:              w,
		expansionsContainer: container.NewVBox(),
		customVarsContainer: container.NewVBox(),
		settingsContainer:   container.NewVBox(),
	}

	state.initUI()
	w.ShowAndRun()
}

func (e *editorState) initUI() {
	// Search and top actions.
	e.searchEntry = widget.NewEntry()
	e.searchEntry.SetPlaceHolder("Search...")
	e.searchEntry.OnChanged = func(string) {
		e.refreshExpansionsView()
	}

	addButton := widget.NewButton("+ Add", func() {
		e.showExpansionForm(nil)
	})
	importButton := widget.NewButton("Import", func() {
		e.importConfig()
	})
	exportButton := widget.NewButton("Export", func() {
		e.exportConfig()
	})

	topBar := container.NewHBox(
		widget.NewLabel("Search:"),
		e.searchEntry,
		layout.NewSpacer(),
		addButton,
		importButton,
		exportButton,
	)

	// Main list of expansions.
	e.refreshExpansionsView()
	expansionsScroll := container.NewVScroll(e.expansionsContainer)
	expansionsScroll.SetMinSize(fyne.NewSize(0, 250))

	// Custom variables and settings.
	e.refreshCustomVars()
	e.refreshSettings()

	bottom := container.NewVBox(
		e.customVarsContainer,
		widget.NewSeparator(),
		e.settingsContainer,
	)

	content := container.NewBorder(
		topBar,
		bottom,
		nil,
		nil,
		expansionsScroll,
	)

	e.window.SetContent(content)
}

func (e *editorState) refreshExpansionsView() {
	all := e.cfg.GetExpansions()
	query := strings.ToLower(strings.TrimSpace(e.searchEntry.Text))

	if query == "" {
		e.filteredExpansions = all
	} else {
		var filtered []config.Expansion
		for _, exp := range all {
			if strings.Contains(strings.ToLower(exp.Trigger), query) ||
				strings.Contains(strings.ToLower(exp.Replacement), query) ||
				strings.Contains(strings.ToLower(exp.Description), query) {
				filtered = append(filtered, exp)
			}
		}
		e.filteredExpansions = filtered
	}

	e.expansionsContainer.Objects = nil

	header := container.NewGridWithColumns(4,
		widget.NewLabelWithStyle("Trigger", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Expansion", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Description", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Actions", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)
	e.expansionsContainer.Add(header)

	for _, exp := range e.filteredExpansions {
		expCopy := exp // capture

		triggerLabel := widget.NewLabel(exp.Trigger)

		preview := exp.Replacement
		if len(preview) > 40 {
			preview = preview[:37] + "..."
		}
		expansionLabel := widget.NewLabel(preview)
		descriptionLabel := widget.NewLabel(exp.Description)

		editBtn := widget.NewButton("Edit", func() {
			e.showExpansionForm(&expCopy)
		})
		deleteBtn := widget.NewButton("Delete", func() {
			e.deleteExpansion(expCopy.Trigger)
		})
		actions := container.NewHBox(editBtn, deleteBtn)

		row := container.NewGridWithColumns(4,
			triggerLabel,
			expansionLabel,
			descriptionLabel,
			actions,
		)

		e.expansionsContainer.Add(row)
	}

	e.expansionsContainer.Refresh()
}

func (e *editorState) showExpansionForm(existing *config.Expansion) {
	var (
		title string
	)

	triggerEntry := widget.NewEntry()
	replacementEntry := widget.NewMultiLineEntry()
	descriptionEntry := widget.NewEntry()
	caseSensitiveCheck := widget.NewCheck("Case sensitive", nil)

	if existing != nil {
		title = "Edit Expansion"
		triggerEntry.SetText(existing.Trigger)
		replacementEntry.SetText(existing.Replacement)
		descriptionEntry.SetText(existing.Description)
		caseSensitiveCheck.SetChecked(existing.CaseSensitive)
	} else {
		title = "Add Expansion"
	}

	formItems := []*widget.FormItem{
		widget.NewFormItem("Trigger", triggerEntry),
		widget.NewFormItem("Replacement", replacementEntry),
		widget.NewFormItem("Description", descriptionEntry),
		widget.NewFormItem("", caseSensitiveCheck),
	}

	onSubmit := func(ok bool) {
		if !ok {
			return
		}

		exp := config.Expansion{
			Trigger:       strings.TrimSpace(triggerEntry.Text),
			Replacement:   replacementEntry.Text,
			CaseSensitive: caseSensitiveCheck.Checked,
			Description:   descriptionEntry.Text,
		}

		if exp.Trigger == "" {
			dialog.ShowError(fmt.Errorf("trigger cannot be empty"), e.window)
			return
		}

		var err error
		if existing == nil {
			err = e.cfg.AddExpansion(exp)
		} else {
			// Remove the old one and then add the updated expansion.
			_ = e.cfg.RemoveExpansion(existing.Trigger)
			err = e.cfg.AddExpansion(exp)
		}
		if err != nil {
			dialog.ShowError(err, e.window)
			return
		}

		if err := e.cfg.Save(); err != nil {
			dialog.ShowError(err, e.window)
			return
		}

		e.refreshExpansionsView()
	}

	dialog.NewForm(title, "Save", "Cancel", formItems, onSubmit, e.window).Show()
}

func (e *editorState) deleteExpansion(trigger string) {
	confirm := dialog.NewConfirm("Delete Expansion",
		fmt.Sprintf("Are you sure you want to delete %q?", trigger),
		func(ok bool) {
			if !ok {
				return
			}
			if err := e.cfg.RemoveExpansion(trigger); err != nil {
				dialog.ShowError(err, e.window)
				return
			}
			if err := e.cfg.Save(); err != nil {
				dialog.ShowError(err, e.window)
				return
			}
			e.refreshExpansionsView()
		}, e.window)
	confirm.Show()
}

func (e *editorState) refreshCustomVars() {
	e.customVarsContainer.Objects = nil

	e.customVarsContainer.Add(widget.NewLabelWithStyle(
		"Custom Variables",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	))

	vars := e.cfg.GetCustomVars()
	if len(vars) == 0 {
		e.customVarsContainer.Add(widget.NewLabel("No custom variables defined."))
	} else {
		keys := make([]string, 0, len(vars))
		for k := range vars {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, key := range keys {
			k := key
			value := vars[k]

			keyLabel := widget.NewLabel(k)
			valueEntry := widget.NewEntry()
			valueEntry.SetText(value)
			valueEntry.OnChanged = func(v string) {
				e.cfg.SetCustomVar(k, v)
				_ = e.cfg.Save()
			}
			removeBtn := widget.NewButton("Remove", func() {
				e.cfg.DeleteCustomVar(k)
				_ = e.cfg.Save()
				e.refreshCustomVars()
			})

			row := container.NewGridWithColumns(3, keyLabel, valueEntry, removeBtn)
			e.customVarsContainer.Add(row)
		}
	}

	addBtn := widget.NewButton("+ Add Variable", func() {
		e.showAddVariableDialog()
	})
	e.customVarsContainer.Add(addBtn)

	e.customVarsContainer.Refresh()
}

func (e *editorState) showAddVariableDialog() {
	keyEntry := widget.NewEntry()
	valueEntry := widget.NewEntry()

	formItems := []*widget.FormItem{
		widget.NewFormItem("Name", keyEntry),
		widget.NewFormItem("Value", valueEntry),
	}

	onSubmit := func(ok bool) {
		if !ok {
			return
		}
		key := strings.TrimSpace(keyEntry.Text)
		if key == "" {
			dialog.ShowError(fmt.Errorf("variable name cannot be empty"), e.window)
			return
		}
		e.cfg.SetCustomVar(key, valueEntry.Text)
		if err := e.cfg.Save(); err != nil {
			dialog.ShowError(err, e.window)
			return
		}
		e.refreshCustomVars()
	}

	dialog.NewForm("Add Custom Variable", "Save", "Cancel", formItems, onSubmit, e.window).Show()
}

func (e *editorState) refreshSettings() {
	e.settingsContainer.Objects = nil

	e.settingsContainer.Add(widget.NewLabelWithStyle(
		"Settings",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	))

	s := e.cfg.GetSettings()

	updateSetting := func(mutator func(*config.Settings)) {
		current := e.cfg.GetSettings()
		mutator(&current)
		e.cfg.UpdateSettings(current)
		_ = e.cfg.Save()
	}

	enabledCheck := widget.NewCheck("Enabled", func(v bool) {
		updateSetting(func(s *config.Settings) { s.Enabled = v })
	})
	enabledCheck.SetChecked(s.Enabled)

	showNotificationsCheck := widget.NewCheck("Show notifications", func(v bool) {
		updateSetting(func(s *config.Settings) { s.ShowNotifications = v })
	})
	showNotificationsCheck.SetChecked(s.ShowNotifications)

	triggerOnSpaceCheck := widget.NewCheck("Trigger on Space", func(v bool) {
		updateSetting(func(s *config.Settings) { s.TriggerOnSpace = v })
	})
	triggerOnSpaceCheck.SetChecked(s.TriggerOnSpace)

	triggerOnTabCheck := widget.NewCheck("Trigger on Tab", func(v bool) {
		updateSetting(func(s *config.Settings) { s.TriggerOnTab = v })
	})
	triggerOnTabCheck.SetChecked(s.TriggerOnTab)

	triggerOnEnterCheck := widget.NewCheck("Trigger on Enter", func(v bool) {
		updateSetting(func(s *config.Settings) { s.TriggerOnEnter = v })
	})
	triggerOnEnterCheck.SetChecked(s.TriggerOnEnter)

	logExpansionsCheck := widget.NewCheck("Log expansions", func(v bool) {
		updateSetting(func(s *config.Settings) { s.LogExpansions = v })
	})
	logExpansionsCheck.SetChecked(s.LogExpansions)

	row1 := container.NewGridWithColumns(2, enabledCheck, showNotificationsCheck)
	row2 := container.NewGridWithColumns(2, triggerOnSpaceCheck, triggerOnTabCheck)
	row3 := container.NewGridWithColumns(2, triggerOnEnterCheck, logExpansionsCheck)

	e.settingsContainer.Add(row1)
	e.settingsContainer.Add(row2)
	e.settingsContainer.Add(row3)

	e.settingsContainer.Refresh()
}

func (e *editorState) importConfig() {
	fd := dialog.NewFileOpen(func(r fyne.URIReadCloser, err error) {
		if err != nil || r == nil {
			return
		}
		defer r.Close()

		data, err := io.ReadAll(r)
		if err != nil {
			dialog.ShowError(err, e.window)
			return
		}

		var imported struct {
			Expansions      []config.Expansion `json:"expansions"`
			CustomVariables map[string]string  `json:"custom_variables"`
			Settings        config.Settings    `json:"settings"`
		}

		if err := json.Unmarshal(data, &imported); err != nil {
			dialog.ShowError(fmt.Errorf("invalid configuration: %w", err), e.window)
			return
		}

		for _, exp := range imported.Expansions {
			_ = e.cfg.AddExpansion(exp)
		}
		for k, v := range imported.CustomVariables {
			e.cfg.SetCustomVar(k, v)
		}

		// We do not blindly overwrite settings; instead, we merge enabled flags
		// and other fields where present.
		current := e.cfg.GetSettings()
		current.Enabled = imported.Settings.Enabled
		current.TriggerOnSpace = imported.Settings.TriggerOnSpace
		current.TriggerOnTab = imported.Settings.TriggerOnTab
		current.TriggerOnEnter = imported.Settings.TriggerOnEnter
		current.ShowNotifications = imported.Settings.ShowNotifications
		current.LogExpansions = imported.Settings.LogExpansions
		e.cfg.UpdateSettings(current)

		if err := e.cfg.Save(); err != nil {
			dialog.ShowError(err, e.window)
			return
		}

		e.refreshExpansionsView()
		e.refreshCustomVars()
		e.refreshSettings()
	}, e.window)

	fd.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
	fd.Show()
}

func (e *editorState) exportConfig() {
	fd := dialog.NewFileSave(func(wc fyne.URIWriteCloser, err error) {
		if err != nil || wc == nil {
			return
		}
		defer wc.Close()

		data := struct {
			Expansions      []config.Expansion `json:"expansions"`
			CustomVariables map[string]string  `json:"custom_variables"`
			Settings        config.Settings    `json:"settings"`
		}{
			Expansions:      e.cfg.GetExpansions(),
			CustomVariables: e.cfg.GetCustomVars(),
			Settings:        e.cfg.GetSettings(),
		}

		bytes, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			dialog.ShowError(err, e.window)
			return
		}

		if _, err := wc.Write(bytes); err != nil {
			dialog.ShowError(err, e.window)
			return
		}
	}, e.window)

	fd.SetFileName("expansions.json")
	fd.Show()
}
