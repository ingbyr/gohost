package form

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gohost/config"
)

var cfg = config.Instance()

type Item interface {
	tea.Model
	Focus(mode FocusMode) tea.Cmd
	Unfocus() tea.Cmd
	SetWidth(width int)
	SetHeight(height int)
	Width() int
	Height() int
	InterceptKey(m tea.KeyMsg) bool
	SetFocusedStyle(style lipgloss.Style)
	SetUnfocusedStyle(style lipgloss.Style)
}
