package form

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gohost/log"
	"gohost/tui/keys"
	"gohost/tui/styles"
)

var _ Item = (*Button)(nil)

func NewButton(text string) *Button {
	return &Button{
		Text:           "[ " + text + " ]",
		OnClick:        func() tea.Cmd { return nil },
		focused:        false,
		focusedStyle:   styles.None,
		unfocusedStyle: styles.None,
		HideFunc:       nil,
	}
}

type Button struct {
	Text           string
	OnClick        func() tea.Cmd
	focused        bool
	focusedStyle   lipgloss.Style
	unfocusedStyle lipgloss.Style
	HideFunc       HideCondition
}

func (b *Button) Focusable() bool {
	return true
}

func (b *Button) Hide() bool {
	if b.HideFunc == nil {
		return false
	}
	return b.HideFunc()
}

func (b *Button) Init() tea.Cmd {
	return nil
}

func (b *Button) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(m, keys.Enter) {
			log.Debug("hit he button by enter")
			return b, b.OnClick()
		}
	}
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
	if key.Matches(m, keys.Enter) {
		return true
	}
	return false
}

func (b *Button) SetFocusedStyle(style lipgloss.Style) {
	b.focusedStyle = style
}

func (b *Button) SetUnfocusedStyle(style lipgloss.Style) {
	b.unfocusedStyle = style
}
