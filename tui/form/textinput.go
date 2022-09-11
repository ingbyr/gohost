package form

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gohost/tui/styles"
)

var _ Item = (*TextInput)(nil)

func NewTextInput() *TextInput {
	t := &TextInput{
		Model:          textinput.New(),
		focusedStyle:   styles.None,
		unfocusedStyle: styles.None,
	}
	t.Unfocus()
	return t
}

type TextInput struct {
	textinput.Model
	focused        bool
	focusedStyle   lipgloss.Style
	unfocusedStyle lipgloss.Style
}

func (t *TextInput) SetFocusedStyle(style lipgloss.Style) {
	t.focusedStyle = style
}

func (t *TextInput) SetUnfocusedStyle(style lipgloss.Style) {
	t.unfocusedStyle = style
}

func (t *TextInput) Init() tea.Cmd {
	return nil
}

func (t *TextInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	t.Model, cmd = t.Model.Update(msg)
	return t, cmd
}

func (t *TextInput) Focus(mode FocusMode) tea.Cmd {
	t.TextStyle = t.focusedStyle
	t.PromptStyle = t.focusedStyle
	t.focused = true
	return t.Model.Focus()
}

func (t *TextInput) Unfocus() tea.Cmd {
	t.TextStyle = t.unfocusedStyle
	t.PromptStyle = t.unfocusedStyle
	t.Model.Blur()
	t.focused = false
	return nil
}

func (t *TextInput) InterceptKey(keyMsg tea.KeyMsg) bool {
	return false
}

func (t *TextInput) View() string {
	if t.focused {
		return t.focusedStyle.Render(t.Model.View())
	}
	return t.unfocusedStyle.Render(t.Model.View())
}
