package view

import (
	tea "github.com/charmbracelet/bubbletea"
	"gohost/tui/widget"
)

type View interface {
	tea.Model
	AddWidget(widget widget.Widget)
	FocusNextWidget() []tea.Cmd
	FocusPreWidget() []tea.Cmd
}
