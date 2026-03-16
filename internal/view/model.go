package view

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Model struct {
	keys       keyMap
	help       help.Model
	inputStyle lipgloss.Style
	lastKey    string
	quitting   bool
	height     int
	width      int
	ready      bool
}

func NewModel() Model {
	return Model{
		keys:       keys,
		help:       help.New(),
		inputStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#FF75B7")),
		height:     10,
		ready:      false,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		m.ready = true
		// If we set a width on the help menu it can gracefully truncate
		// its view as needed.
		m.help.SetWidth(msg.Width)

	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() tea.View {
	if m.quitting {
		view := tea.NewView("Goodbye!\n")
		view.AltScreen = true
		return view
	}

	if !m.ready {
		view := tea.NewView("booting...\n")
		view.AltScreen = true
		return view
	}

	helpView := m.help.View(m.keys)
	// 3 = one line for index, one line for bottom padding and one line for line between the help and the main content
	availableHeight := m.height - 3 - strings.Count(helpView, "\n")

	s := ""
	for i := range availableHeight {
		s += "["
		s += fmt.Sprintf("%d  %d", availableHeight, i+1)
		s += "\n"
	}

	view := tea.NewView(s + "\n" + helpView)
	view.AltScreen = true
	return view
}
