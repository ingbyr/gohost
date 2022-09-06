package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type TextView struct {
	model        *Model
	HostTextarea textarea.Model
}

func NewTextView(model *Model) *TextView {
	t := textarea.New()
	t.Placeholder = "Host items here"
	t.Focus()
	return &TextView{
		model:        model,
		HostTextarea: textarea.New(),
	}
}

func (v *TextView) Init() tea.Cmd {
	return nil
}

func (v *TextView) Update(msg tea.Msg) []tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.HostTextarea.SetHeight(msg.Height - 10)
		v.HostTextarea.SetWidth(msg.Width - v.model.groupView.groupList.Width())
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Enter):
		}
	}
	v.HostTextarea, cmd = v.HostTextarea.Update(msg)
	return append(cmds, cmd)
}

func (v *TextView) View() string {
	return v.HostTextarea.View()
}
