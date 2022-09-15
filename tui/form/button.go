package form

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"gohost/log"
	"gohost/tui/keys"
)

var _ ItemModel = (*Button)(nil)

func NewButton(text string) *Button {
	return &Button{
		CommonItem: NewCommonItem(),
		Text:       "[ " + text + " ]",
		OnClick:    func() tea.Cmd { return nil },
	}
}

type Button struct {
	*CommonItem
	Text    string
	OnClick func() tea.Cmd
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

func (b *Button) InterceptKey(m tea.KeyMsg) bool {
	if key.Matches(m, keys.Enter) {
		return true
	}
	return false
}
