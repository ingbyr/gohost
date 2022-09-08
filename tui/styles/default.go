package styles

import "github.com/charmbracelet/lipgloss"

const padding = 0

var (
	None = lipgloss.NewStyle()
	DefaultView = None
	FocusedView = None
	FocusedWidget = None

	//DefaultView = lipgloss.NewStyle().
	//		PaddingLeft(padding).
	//		PaddingRight(padding).
	//		BorderStyle(lipgloss.NormalBorder()).
	//		BorderForeground(lipgloss.Color("70"))
	//
	//FocusedView = lipgloss.NewStyle().
	//		PaddingLeft(padding).
	//		PaddingRight(padding).
	//		BorderStyle(lipgloss.NormalBorder()).
	//		BorderForeground(lipgloss.Color("69"))
	//
	//FocusedWidget = lipgloss.NewStyle().
	//		Foreground(lipgloss.Color("69"))
)
