package widget

import tea "github.com/charmbracelet/bubbletea"

type Widget interface {
	tea.Model
	Focus() tea.Cmd
	Unfocus() tea.Cmd
	HasNext() bool
	HasPre() bool
}
