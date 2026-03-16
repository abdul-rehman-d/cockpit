package view

import (
	"fmt"
	"os"
	osuser "os/user"
	"strings"
	"time"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/progress"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/abdul-rehman-d/cockpit/internal/usage"
	"github.com/abdul-rehman-d/cockpit/internal/utils"
)

type tickMsg time.Time

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
	now      time.Time

	// reusable progress bar
	prog progress.Model

	// providers
	service *usage.Service
	user    string

	// misc ui stuff
	maxSampleNameLength int
	textStyle           lipgloss.Style
}

func NewModel() Model {
	service := usage.NewService()
	helpModel := help.New()
	helpModel.Styles = utils.DraculaHelpStyles()

	samples := service.GetAllSamples()
	maxSampleNameLength := 0
	for _, sample := range samples {
		maxSampleNameLength = max(maxSampleNameLength, len(sample.Name))
	}

	user := currentUserName()

	return Model{
		keys: keys,
		help: helpModel,

		prog: utils.NewDraculaProgress(),

		service: service,
		user:    user,

		maxSampleNameLength: maxSampleNameLength,
		textStyle:           utils.NormalTextStyle,
		now:                 time.Now(),
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
		m.now = time.Time(msg)
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

	s := "\n" + m.renderGreetings() + "\n\n"
	s += m.renderClock() + "\n\n"

	samples := m.service.GetAllSamples()
	for _, sample := range samples {
		s += m.renderSample(sample)
	}

	availableHeight -= strings.Count(s, "\n")
	s += strings.Repeat("\n", availableHeight)

	s += "\n" + helpView

	mainWindow := lipgloss.NewStyle().Width(m.width)

	view := tea.NewView(mainWindow.Render(s))
	view.AltScreen = true
	return view
}

func (m Model) renderSample(sample usage.Sample) string {
	row1 := m.textStyle.Render(sample.Name+
		strings.Repeat(" ", m.maxSampleNameLength-len(sample.Name)+1)) +
		m.prog.ViewAs(sample.Value)
	row2 := strings.Repeat(" ", m.maxSampleNameLength+1) +
		m.textStyle.Render(sample.ValueInWords)

	return row1 + "\n" + row2 + "\n\n"
}

func (m Model) renderClock() string {
	hour := m.now.Hour() % 12
	if hour == 0 {
		hour = 12
	}
	separator := ":"
	if m.now.Second()%2 == 1 {
		separator = " "
	}
	meridiem := "AM"
	if m.now.Hour() >= 12 {
		meridiem = "PM"
	}

	clock := fmt.Sprintf("%02d%s%02d (%s)", hour, separator, m.now.Minute(), meridiem)
	return m.textStyle.Render(clock)
}

func (m Model) renderGreetings() string {
	if m.user == "" {
		return m.textStyle.Render("greetings")
	}
	return m.textStyle.Render(fmt.Sprintf("greetings, %s", m.user))
}

func currentUserName() string {
	if u, err := osuser.Current(); err == nil {
		if u.Username != "" {
			return u.Username
		}
		if u.Name != "" {
			return u.Name
		}
	}

	if username := os.Getenv("USER"); username != "" {
		return username
	}
	if username := os.Getenv("USERNAME"); username != "" {
		return username
	}

	return ""
}

func doTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
