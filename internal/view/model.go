package view

import (
	"strings"
	"time"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/progress"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/abdul-rehman-d/cockpit/internal/usage"
)

type tickMsg time.Time

var (
	yellow = lipgloss.Color("#FDFF8C")
	pink   = lipgloss.Color("#FF7CCB")
)

type Model struct {
	// help menu thingy
	keys keyMap
	help help.Model

	// dimensions of the terminal
	width  int
	height int

	// state
	quitting bool
	ready    bool

	// reusable progress bar
	prog progress.Model

	// usage service
	service *usage.Service

	// misc ui stuff
	maxSampleNameLength int
}

func NewModel() Model {
	service := usage.NewService()
	samples := service.GetAllSamples()
	maxSampleNameLength := 0
	for _, sample := range samples {
		maxSampleNameLength = max(maxSampleNameLength, len(sample.Name))
	}
	return Model{
		keys: keys,
		help: help.New(),

		prog: progress.New(progress.WithScaled(true), progress.WithColors(pink, yellow)),

		service: service,

		maxSampleNameLength: maxSampleNameLength,
	}
}

func (m Model) Init() tea.Cmd {
	return doTick()
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
		m.prog.SetWidth(msg.Width - m.maxSampleNameLength - 5) // 5 = one space between name and bar and 4 spaces of padding

	case tickMsg:
		return m, doTick()

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

	samples := m.service.GetAllSamples()
	for _, sample := range samples {
		s += m.renderSample(sample)
	}

	availableHeight -= strings.Count(s, "\n")
	s += strings.Repeat("\n", availableHeight)

	view := tea.NewView(s + "\n" + helpView)
	view.AltScreen = true
	return view
}

func (m Model) renderSample(sample usage.Sample) string {
	row1 := sample.Name +
		strings.Repeat(" ", m.maxSampleNameLength-len(sample.Name)+1) +
		m.prog.ViewAs(sample.Value)
	row2 := strings.Repeat(" ", m.maxSampleNameLength+1) +
		sample.ValueInWords

	return row1 + "\n" + row2 + "\n\n"
}

func doTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
