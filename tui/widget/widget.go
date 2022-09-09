package widget

import tea "github.com/charmbracelet/bubbletea"

type Widget interface {
	tea.Model
	Focus() tea.Cmd
	Unfocus() tea.Cmd
	HandleKeyUp() bool
	HandleKeyDown() bool
	SetWidth(width int)
	SetHeight(height int) int
}
