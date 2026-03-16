package utils

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/progress"
	"charm.land/lipgloss/v2"
)

var (
	DraculaBackground = lipgloss.Color("#282A36")
	DraculaForeground = lipgloss.Color("#F8F8F2")
	DraculaComment    = lipgloss.Color("#6272A4")
	DraculaPurple     = lipgloss.Color("#BD93F9")
	DraculaPink       = lipgloss.Color("#FF79C6")
)

var NormalTextStyle = lipgloss.NewStyle().Foreground(DraculaForeground)

func NewDraculaProgress() progress.Model {
	p := progress.New(
		progress.WithScaled(true),
		progress.WithColors(DraculaPurple, DraculaPink),
	)
	p.EmptyColor = DraculaComment
	p.PercentageStyle = NormalTextStyle.Background(DraculaBackground)
	p.PercentFormat = " %05.2f%%"
	return p
}

func DraculaHelpStyles() help.Styles {
	styles := help.DefaultDarkStyles()
	commentStyle := lipgloss.NewStyle().Background(DraculaBackground).Foreground(DraculaComment)

	styles.Ellipsis = commentStyle
	styles.ShortKey = commentStyle
	styles.ShortDesc = commentStyle
	styles.ShortSeparator = commentStyle
	styles.FullKey = commentStyle
	styles.FullDesc = commentStyle
	styles.FullSeparator = commentStyle

	return styles
}
