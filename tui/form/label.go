package form

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var _ Item = (*Label)(nil)

func NewLabel(label string) *Label {
	return &Label{label: label}
}

type Label struct {
	label         string
	width, height int
}

func (l *Label) Focusable() bool {
	return false
}

func (l *Label) Init() tea.Cmd {
	return nil
}

func (l *Label) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		l.width, l.height = msg.Width, msg.Height
	}
	return l, nil
}

func (l *Label) View() string {
	return l.label
}

func (l *Label) Focus(mode FocusMode) tea.Cmd {
	return nil
}

func (l *Label) Unfocus() tea.Cmd {
	return nil
}

func (l *Label) InterceptKey(m tea.KeyMsg) bool {
	return false
}

func (l *Label) SetFocusedStyle(style lipgloss.Style) {
}

func (l *Label) SetUnfocusedStyle(style lipgloss.Style) {
}

func (l *Label) Hide() bool {
	return false
}
