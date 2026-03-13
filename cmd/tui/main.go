package main

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const tickRate = time.Second

var (
	draculaBg      = lipgloss.Color("#282A36")
	draculaFg      = lipgloss.Color("#F8F8F2")
	draculaCyan    = lipgloss.Color("#8BE9FD")
	draculaGreen   = lipgloss.Color("#50FA7B")
	draculaComment = lipgloss.Color("#6272A4")
	draculaPink    = lipgloss.Color("#FF79C6")
)

type tickMsg time.Time

type model struct {
	width    int
	height   int
	start    time.Time
	lastTick time.Time
	ticks    int
	realFPS  float64
	now      time.Time
	blinkOn  bool
	ready    bool
}

func newModel() model {
	now := time.Now()
	return model{
		start:    now,
		lastTick: now,
		now:      now,
		blinkOn:  true,
	}
}

func nextTick() tea.Cmd {
	return tea.Tick(tickRate, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	return nextTick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true

	case tickMsg:
		now := time.Time(msg)
		delta := now.Sub(m.lastTick)
		if delta > 0 {
			m.realFPS = 1 / delta.Seconds()
		}
		m.lastTick = now
		m.now = now
		m.blinkOn = !m.blinkOn
		m.ticks++
		return m, nextTick()
	}

	return m, nil
}

func (m model) View() string {
	if !m.ready {
		return "Booting cockpit..."
	}

	innerWidth := max(1, m.width-2)
	innerHeight := max(1, m.height-2)

	uptime := m.now.Sub(m.start).Round(time.Second)
	statsText := fmt.Sprintf("fps: %4.1f  uptime: %s", m.realFPS, uptime)

	colon := ":"
	if !m.blinkOn {
		colon = " "
	}
	clockText := fmt.Sprintf("%s%s%s", m.now.Format("15"), colon, m.now.Format("04"))

	statsLine := lipgloss.NewStyle().Foreground(draculaGreen).Render(statsText)
	clockLine := lipgloss.NewStyle().Bold(true).Foreground(draculaPink).Render(clockText)
	footerLine := lipgloss.NewStyle().Foreground(draculaComment).Render("Press q / esc to quit")

	lines := make([]string, innerHeight)
	for i := range lines {
		lines[i] = strings.Repeat(" ", innerWidth)
	}

	lines[0] = padOrTrim(statsLine, innerWidth)

	clockRow := innerHeight / 2
	lines[clockRow] = centerText(clockLine, innerWidth)

	lines[innerHeight-1] = padOrTrim(footerLine, innerWidth)

	content := strings.Join(lines, "\n")

	panel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(draculaCyan).
		Foreground(draculaFg).
		Background(draculaBg).
		Width(innerWidth).
		Height(innerHeight).
		Render(content)

	return lipgloss.NewStyle().Background(draculaBg).Render(panel)
}

func padOrTrim(s string, width int) string {
	if width <= 0 {
		return ""
	}
	w := lipgloss.Width(s)
	if w >= width {
		return lipgloss.NewStyle().MaxWidth(width).Render(s)
	}
	return s + strings.Repeat(" ", width-w)
}

func centerText(s string, width int) string {
	if width <= 0 {
		return ""
	}
	w := lipgloss.Width(s)
	if w >= width {
		return lipgloss.NewStyle().MaxWidth(width).Render(s)
	}
	left := (width - w) / 2
	right := width - w - left
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
}

func main() {
	p := tea.NewProgram(newModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("error: %v\n", err)
	}
}
