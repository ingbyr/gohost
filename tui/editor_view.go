package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"gohost/gohost"
)

type EditorView struct {
	model        *Model
	HostTextarea textarea.Model
	host         gohost.Host
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
		v.HostTextarea.SetHeight(msg.Height - v.model.reservedHeight)
		v.HostTextarea.SetWidth(msg.Width - v.model.groupView.groupList.Width())
	}
	if v.model.state == editorViewState {
		v.HostTextarea, cmd = v.HostTextarea.Update(msg)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.Save):
				v.host.SetContent([]byte(v.HostTextarea.Value()))
				err := gohost.GetService().UpdateHost(v.host)
				if err != nil {
					v.model.Log(err.Error())
				}
			}
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

func (v *EditorView) SetHost(host gohost.Host) {
	v.host = host
	v.HostTextarea.SetValue(string(host.GetContent()))
}
