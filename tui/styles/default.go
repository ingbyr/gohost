package styles

import "github.com/charmbracelet/lipgloss"

var (
	DefaultView  = lipgloss.NewStyle().PaddingLeft(1).PaddingRight(1).BorderStyle(lipgloss.HiddenBorder())
	FocusedView  = lipgloss.NewStyle().PaddingLeft(1).PaddingRight(1).BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("69"))
	FocusedModel = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	None         = lipgloss.NewStyle()
)
