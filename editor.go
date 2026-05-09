package vtea

import tea "charm.land/bubbletea/v2"

// Editor implements tea.Model by embedding *Model
type Editor struct {
	*Model
}

// NewEditor creates a new editor instance with the provided options
func NewEditor(opts ...EditorOption) Editor {
	return Editor{Model: New(opts...)}
}

// Update handles messages and updates the editor state
// This is part of the tea.Model interface
func (e Editor) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m, cmd := e.Model.Update(msg)
	e.Model = m
	return e, cmd
}

func (e Editor) View() tea.View {
	view := tea.NewView(e.Model.View())
	if e.Model.fullScreen {
		view.AltScreen = true
	}
	return view
}
