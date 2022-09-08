package styles

import "github.com/charmbracelet/lipgloss"

const padding = 3

var (
	DefaultView   = lipgloss.NewStyle().PaddingLeft(padding).PaddingRight(padding).BorderStyle(lipgloss.HiddenBorder())
	WidgetView    = lipgloss.NewStyle().PaddingBottom(padding).BorderStyle(lipgloss.HiddenBorder())
	FocusedView   = lipgloss.NewStyle().PaddingLeft(padding).PaddingRight(padding).BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("69"))
	FocusedWidget = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	None          = lipgloss.NewStyle()
)
