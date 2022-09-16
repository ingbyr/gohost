package styles

import "github.com/charmbracelet/lipgloss"

const padding = 1

var (
	None        = lipgloss.NewStyle()
	DefaultView = lipgloss.NewStyle().
			PaddingLeft(padding).
			PaddingRight(padding).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.AdaptiveColor{Light: "244", Dark: "244"})
	FocusedView = lipgloss.NewStyle().
			PaddingLeft(padding).
			PaddingRight(padding).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.AdaptiveColor{Light: "28", Dark: "34"})
	UnfocusedFormItem = lipgloss.NewStyle()
	FocusedFormItem   = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "28", Dark: "34"})

	StatusLine = lipgloss.NewStyle()
)
