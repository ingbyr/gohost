package form

import (
	tea "github.com/charmbracelet/bubbletea"
	"gohost/config"
)

var cfg = config.Instance()

type Widget interface {
	tea.Model
	Focus(mode FocusMode) tea.Cmd
	Unfocus() tea.Cmd
	HandleKeyUp() bool
	HandleKeyDown() bool
	SetWidth(width int)
	SetHeight(height int)
	Width() int
	Height() int
}
