package utils

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/progress"
	"charm.land/lipgloss/v2"
)

var (
	// ANSI 16-color palette (4-bit) for broad TTY support.
	DraculaForeground = lipgloss.Color("7")
	DraculaComment    = lipgloss.Color("8")
	DraculaPurple     = lipgloss.Color("5")
	DraculaPink       = lipgloss.Color("13")
)

var NormalTextStyle = lipgloss.NewStyle().Foreground(DraculaForeground)

func NewDraculaProgress() progress.Model {
	p := progress.New(
		progress.WithScaled(true),
		progress.WithFillCharacters('█', ' '),
		progress.WithColors(DraculaPurple),
	)
	p.PercentageStyle = NormalTextStyle
	p.PercentFormat = " %05.2f%%"
	return p
}

func DraculaHelpStyles() help.Styles {
	styles := help.DefaultDarkStyles()
	commentStyle := lipgloss.NewStyle().Foreground(DraculaComment)

	styles.Ellipsis = commentStyle
	styles.ShortKey = commentStyle
	styles.ShortDesc = commentStyle
	styles.ShortSeparator = commentStyle
	styles.FullKey = commentStyle
	styles.FullDesc = commentStyle
	styles.FullSeparator = commentStyle

	return styles
}
