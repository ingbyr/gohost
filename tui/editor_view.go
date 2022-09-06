package tui

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type EditorView struct {
	model        *Model
	HostTextarea textarea.Model
}

func NewTextView(model *Model) *EditorView {
	t := textarea.New()
	t.ShowLineNumbers = true
	t.Placeholder = "Host items here"
	return &EditorView{
		model:        model,
		HostTextarea: textarea.New(),
	}
}

func (v *EditorView) Init() tea.Cmd {
	return nil
}

func (v *EditorView) Update(msg tea.Msg) []tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.HostTextarea.SetHeight(msg.Height - v.model.helpView.MaxHeight())
		v.HostTextarea.SetWidth(msg.Width - v.model.groupView.groupList.Width())
	case tea.KeyMsg:
		if v.model.state == editorViewState {
			v.HostTextarea, cmd = v.HostTextarea.Update(msg)
		}
	}
	return append(cmds, cmd)
}

func (v *EditorView) View() string {
	return v.HostTextarea.View()
}

func (v *EditorView) Focus() {
	v.HostTextarea.Focus()
}
