package widget

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"gohost/tui/styles"
)

var _ Widget = (*TextInput)(nil)

func NewTextInput() *TextInput {
	t := &TextInput{
		Model: textinput.New(),
	}
	t.Unfocus()
	return t
}

type TextInput struct {
	textinput.Model
}

func (t *TextInput) Width() int {
	return t.Model.Width
}

func (t *TextInput) Height() int {
	return 1
}

func (t *TextInput) SetWidth(width int) {
	t.Model.Width = width - len(t.Prompt) - 1
}

func (t *TextInput) SetHeight(height int) {
	return
}

func (t *TextInput) Init() tea.Cmd {
	return nil
}

func (t *TextInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	t.Model, cmd = t.Model.Update(msg)
	return t, cmd
}

func (t *TextInput) Focus() tea.Cmd {
	t.TextStyle = styles.FocusedWidget
	t.PromptStyle = styles.FocusedWidget
	return t.Model.Focus()
}

func (t *TextInput) Unfocus() tea.Cmd {
	t.TextStyle = styles.None
	t.PromptStyle = styles.None
	t.Model.Blur()
	return nil
}

func (t *TextInput) HandleKeyUp() bool {
	return false
}

func (t *TextInput) HandleKeyDown() bool {
	return false
}
