package gui

import (
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"text-expander/config"
)

type editorState struct {
	cfg                 *config.Config
	window              fyne.Window
	searchEntry         *widget.Entry
	categoryFilter      *widget.Select
	expansionsContainer *fyne.Container
	customVarsContainer *fyne.Container
	settingsContainer   *fyne.Container
	filteredExpansions  []config.Expansion
}

// CreateEditorWindow creates the editor UI in the given window
func CreateEditorWindow(w fyne.Window, cfg *config.Config) {
	state := &editorState{
		cfg:                cfg,
		window:             w,
		filteredExpansions: cfg.GetExpansions(),
	}

	state.initUI()
}

// ShowEditor launches the enhanced configuration editor (kept for backwards compatibility)
// Now just launches separate process
func ShowEditor(a fyne.App, cfg *config.Config) {
	// This function is now deprecated - main.go should launch gui-config.exe instead
	// Keeping it here to avoid breaking existing code
	w := a.NewWindow("Text Expander Manager")
	w.Resize(fyne.NewSize(1000, 700))
	w.CenterOnScreen()

	CreateEditorWindow(w, cfg)
	w.Show()
}

func (s *editorState) initUI() {
	// Initialize filtered expansions first
	s.filteredExpansions = s.cfg.GetExpansions()

	// Create search bar
	s.searchEntry = widget.NewEntry()
	s.searchEntry.SetPlaceHolder("🔍 Search expansions...")
	s.searchEntry.OnChanged = func(query string) {
		s.filterExpansions(query)
	}

	// Create category filter
	categories := []string{"All", "Python", "JavaScript", "Go", "C", "SQL", "HTML", "CSS", "Personal", "Professional", "Symbols"}
	s.categoryFilter = widget.NewSelect(categories, func(selected string) {
		s.filterExpansions(s.searchEntry.Text)
	})
	s.categoryFilter.PlaceHolder = "Category"
	// Don't set selected yet - will do after creating containers

	// Create new expansion button
	newBtn := widget.NewButton("+ New Expansion", func() {
		ShowExpansionDialog(s.window, s.cfg, nil, func() {
			s.refreshExpansionsView()
		})
	})
	newBtn.Importance = widget.HighImportance

	// Help button
	helpBtn := widget.NewButton("?", func() {
		ShowHelpDialog(s.window)
	})
	helpBtn.Importance = widget.LowImportance

	// Top toolbar with generous spacing
	toolbar := container.NewBorder(
		nil, nil,
		container.NewPadded(container.NewHBox(layout.NewSpacer())),
		container.NewPadded(container.NewHBox(newBtn, helpBtn)),
		container.NewPadded(s.searchEntry), // Search bar fills remaining space
	)

	// Create tab container
	tabs := container.NewAppTabs(
		container.NewTabItem("Expansions", s.createExpansionsTab()),
		container.NewTabItem("Variables", s.createVariablesTab()),
		container.NewTabItem("Settings", s.createSettingsTab()),
	)

	// Main layout
	content := container.NewBorder(toolbar, nil, nil, nil, tabs)
	s.window.SetContent(content)

	// Now safe to set default category (after containers are created)
	s.categoryFilter.SetSelected("All")
}

func (s *editorState) createExpansionsTab() *fyne.Container {
	s.expansionsContainer = container.NewVBox()
	s.refreshExpansionsView()

	scroll := container.NewVScroll(s.expansionsContainer)
	scroll.SetMinSize(fyne.NewSize(900, 500))

	// Add padding around the scroll area for breathing room
	return container.NewPadded(scroll)
}

func (s *editorState) createVariablesTab() *fyne.Container {
	s.customVarsContainer = container.NewVBox()
	s.refreshCustomVars()

	addVarBtn := widget.NewButton("+ Add Variable", func() {
		s.showAddVariableDialog()
	})

	scroll := container.NewVScroll(s.customVarsContainer)
	return container.NewBorder(addVarBtn, nil, nil, nil, scroll)
}

func (s *editorState) createSettingsTab() *fyne.Container {
	s.settingsContainer = container.NewVBox()
	s.refreshSettings()

	scroll := container.NewVScroll(s.settingsContainer)
	return container.NewBorder(nil, nil, nil, nil, scroll)
}

func (s *editorState) refreshExpansionsView() {
	s.expansionsContainer.Objects = nil

	if len(s.filteredExpansions) == 0 {
		emptyLabel := widget.NewLabel("No expansions found. Click '+ New Expansion' to add one!")
		s.expansionsContainer.Add(emptyLabel)
		s.expansionsContainer.Refresh()
		return
	}

	// Create expansion cards
	for i := range s.filteredExpansions {
		exp := &s.filteredExpansions[i]

		card := NewExpansionCard(
			exp,
			func(e *config.Expansion) {
				ShowExpansionDialog(s.window, s.cfg, e, func() {
					s.refreshExpansionsView()
				})
			},
			func(trigger string) {
				ShowDeleteConfirmation(s.window, trigger, func() {
					s.deleteExpansion(trigger)
				})
			},
		)

		s.expansionsContainer.Add(card)
	}

	s.expansionsContainer.Refresh()
}

func (s *editorState) filterExpansions(query string) {
	allExpansions := s.cfg.GetExpansions()
	s.filteredExpansions = []config.Expansion{}

	query = strings.ToLower(query)
	selectedCategory := s.categoryFilter.Selected

	for _, exp := range allExpansions {
		// Category filter - handle empty categories
		if selectedCategory != "" && selectedCategory != "All" {
			// If expansion has no category, skip it unless "All" is selected
			if exp.Category == "" {
				continue // Skip uncategorized expansions when filtering by category
			}
			// Check if category matches
			if exp.Category != selectedCategory {
				continue
			}
		}

		// Text search
		if query != "" {
			match := strings.Contains(strings.ToLower(exp.Trigger), query) ||
				strings.Contains(strings.ToLower(exp.Description), query) ||
				strings.Contains(strings.ToLower(exp.Replacement), query)

			if !match {
				continue
			}
		}

		s.filteredExpansions = append(s.filteredExpansions, exp)
	}

	// Sort alphabetically by trigger
	sort.Slice(s.filteredExpansions, func(i, j int) bool {
		return s.filteredExpansions[i].Trigger < s.filteredExpansions[j].Trigger
	})

	s.refreshExpansionsView()
}

func (s *editorState) deleteExpansion(trigger string) {
	if err := s.cfg.RemoveExpansion(trigger); err != nil {
		dialog.ShowError(err, s.window)
		return
	}

	if err := s.cfg.Save(); err != nil {
		dialog.ShowError(err, s.window)
		return
	}

	s.filteredExpansions = s.cfg.GetExpansions()
	s.refreshExpansionsView()
}

func (s *editorState) refreshCustomVars() {
	s.customVarsContainer.Objects = nil

	vars := s.cfg.GetCustomVars()
	if len(vars) == 0 {
		s.customVarsContainer.Add(widget.NewLabel("No custom variables defined"))
		s.customVarsContainer.Refresh()
		return
	}

	for key, value := range vars {
		varCard := container.NewHBox(
			widget.NewLabel(key),
			widget.NewLabel("="),
			widget.NewLabel(value),
			layout.NewSpacer(),
			widget.NewButton("Delete", func(k string) func() {
				return func() {
					s.cfg.DeleteCustomVar(k)
					s.cfg.Save()
					s.refreshCustomVars()
				}
			}(key)),
		)
		s.customVarsContainer.Add(varCard)
	}

	s.customVarsContainer.Refresh()
}

func (s *editorState) showAddVariableDialog() {
	keyEntry := widget.NewEntry()
	keyEntry.SetPlaceHolder("Variable name (e.g., NAME)")

	valueEntry := widget.NewEntry()
	valueEntry.SetPlaceHolder("Variable value")

	form := container.NewVBox(
		widget.NewLabel("Variable Name:"),
		keyEntry,
		widget.NewLabel("Value:"),
		valueEntry,
	)

	dialog.NewCustomConfirm("Add Custom Variable", "Add", "Cancel", form, func(add bool) {
		if add && keyEntry.Text != "" && valueEntry.Text != "" {
			s.cfg.SetCustomVar(keyEntry.Text, valueEntry.Text)
			s.cfg.Save()
			s.refreshCustomVars()
		}
	}, s.window).Show()
}

func (s *editorState) refreshSettings() {
	s.settingsContainer.Objects = nil

	settings := s.cfg.GetSettings()

	// Create settings controls
	enabledCheck := widget.NewCheck("Enable expansions", func(checked bool) {
		settings.Enabled = checked
		s.cfg.UpdateSettings(settings)
		s.cfg.Save()
	})
	enabledCheck.SetChecked(settings.Enabled)

	spaceCheck := widget.NewCheck("Trigger on Space", func(checked bool) {
		settings.TriggerOnSpace = checked
		s.cfg.UpdateSettings(settings)
		s.cfg.Save()
	})
	spaceCheck.SetChecked(settings.TriggerOnSpace)

	tabCheck := widget.NewCheck("Trigger on Tab", func(checked bool) {
		settings.TriggerOnTab = checked
		s.cfg.UpdateSettings(settings)
		s.cfg.Save()
	})
	tabCheck.SetChecked(settings.TriggerOnTab)

	enterCheck := widget.NewCheck("Trigger on Enter", func(checked bool) {
		settings.TriggerOnEnter = checked
		s.cfg.UpdateSettings(settings)
		s.cfg.Save()
	})
	enterCheck.SetChecked(settings.TriggerOnEnter)

	notificationsCheck := widget.NewCheck("Show notifications", func(checked bool) {
		settings.ShowNotifications = checked
		s.cfg.UpdateSettings(settings)
		s.cfg.Save()
	})
	notificationsCheck.SetChecked(settings.ShowNotifications)

	loggingCheck := widget.NewCheck("Log expansions", func(checked bool) {
		settings.LogExpansions = checked
		s.cfg.UpdateSettings(settings)
		s.cfg.Save()
	})
	loggingCheck.SetChecked(settings.LogExpansions)

	s.settingsContainer.Add(widget.NewLabelWithStyle("General", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
	s.settingsContainer.Add(enabledCheck)
	s.settingsContainer.Add(widget.NewSeparator())

	s.settingsContainer.Add(widget.NewLabelWithStyle("Trigger Keys", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
	s.settingsContainer.Add(spaceCheck)
	s.settingsContainer.Add(tabCheck)
	s.settingsContainer.Add(enterCheck)
	s.settingsContainer.Add(widget.NewSeparator())

	s.settingsContainer.Add(widget.NewLabelWithStyle("Visual Feedback", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
	s.settingsContainer.Add(notificationsCheck)
	s.settingsContainer.Add(widget.NewSeparator())

	s.settingsContainer.Add(widget.NewLabelWithStyle("Logging", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
	s.settingsContainer.Add(loggingCheck)

	s.settingsContainer.Refresh()
}
