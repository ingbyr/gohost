package form

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var _ ItemModel = (*TextInput)(nil)

func NewTextInput() *TextInput {
	t := &TextInput{
		CommonItem: NewCommonItem(),
		Model:      textinput.New(),
	}
	t.Unfocus()
	return t
}

type TextInput struct {
	*CommonItem
	textinput.Model
}

func (t *TextInput) Init() tea.Cmd {
	return nil
}

func (t *TextInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	t.Model, cmd = t.Model.Update(msg)
	return t, cmd
}

func (t *TextInput) View() string {
	if t.focused {
		return t.focusedStyle.Render(t.Model.View())
	}
	return t.unfocusedStyle.Render(t.Model.View())
}

func (t *TextInput) Focus(mode FocusMode) tea.Cmd {
	t.CommonItem.Focus(mode)
	t.Model.TextStyle = t.focusedStyle
	t.Model.PromptStyle = t.focusedStyle
	return t.Model.Focus()
}

func (t *TextInput) Unfocus() tea.Cmd {
	t.CommonItem.Unfocus()
	t.Model.TextStyle = t.unfocusedStyle
	t.Model.PromptStyle = t.unfocusedStyle
	t.Model.Blur()
	return nil
}
