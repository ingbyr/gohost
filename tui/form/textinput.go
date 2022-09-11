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
		width:          0,
		height:         1,
		focusedStyle:   styles.None,
		unfocusedStyle: styles.None,
	}
	t.Unfocus()
	return t
}

type TextInput struct {
	textinput.Model
	width, height  int
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

func (t *TextInput) Width() int {
	return t.Model.Width
}

func (t *TextInput) Height() int {
	if t.focused {
		return 1 + t.focusedStyle.GetHeight()
	}
	return 1 + t.unfocusedStyle.GetHeight()
}

func (t *TextInput) SetWidth(width int) {
	t.Model.Width = width - len(t.Prompt) - 1
	t.width = width
}

func (t *TextInput) SetHeight(height int) {
	if height > 0 {
		t.height = 1
	} else {
		t.height = 0
	}
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
	if t.height <= 0 {
		return ""
	}
	if t.focused {
		return t.focusedStyle.Render(t.Model.View())
	}
	return t.unfocusedStyle.Render(t.Model.View())
}
