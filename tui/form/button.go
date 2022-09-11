package form

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gohost/tui/styles"
)

func NewButton(text string) *Button {
	return &Button{
		Text:           text,
		focused:        false,
		focusedStyle:   styles.None,
		unfocusedStyle: styles.None,
	}
}

type Button struct {
	Text           string
	focused        bool
	focusedStyle   lipgloss.Style
	unfocusedStyle lipgloss.Style
}

func (b *Button) Init() tea.Cmd {
	return nil
}

func (b *Button) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return b, nil
}

func (b *Button) View() string {
	if b.focused {
		return b.focusedStyle.Render(b.Text)
	}
	return b.unfocusedStyle.Render(b.Text)
}

func (b *Button) Focus(mode FocusMode) tea.Cmd {
	b.focused = true
	return nil
}

func (b *Button) Unfocus() tea.Cmd {
	b.focused = false
	return nil
}

func (b *Button) InterceptKey(m tea.KeyMsg) bool {
	return false
}

func (b *Button) SetFocusedStyle(style lipgloss.Style) {
	b.focusedStyle = style
}

func (b *Button) SetUnfocusedStyle(style lipgloss.Style) {
	b.unfocusedStyle = style
}
