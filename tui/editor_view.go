package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"gohost/config"
	"gohost/gohost"
	"os"
)

type EditorView struct {
	model      *Model
	hostEditor textarea.Model
	host       gohost.Host
}

func NewTextView(model *Model) *EditorView {
	hostEditor := textarea.New()
	hostEditor.ShowLineNumbers = true
	return &EditorView{
		model:      model,
		hostEditor: hostEditor,
	}
}

func (v *EditorView) Init() tea.Cmd {
	return func() tea.Msg {
		sysHost, err := os.ReadFile(config.Instance().SysHostFile)
		if err != nil {
			v.hostEditor.SetValue("Can not open system hosts file")
			return nil
		}
		v.hostEditor.SetValue(string(sysHost))
		return nil
	}
}

func (v *EditorView) Update(msg tea.Msg) []tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	v.hostEditor, cmd = v.hostEditor.Update(msg)
	cmds = append(cmds, cmd)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.hostEditor.SetHeight(msg.Height - v.model.reservedHeight)
		v.hostEditor.SetWidth(msg.Width - v.model.groupView.groupList.Width())
	}
	if v.model.state == editorViewState {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.Save):
				v.host.SetContent([]byte(v.hostEditor.Value()))
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
	return v.hostEditor.View()
}

func (v *EditorView) Focus() {
	v.hostEditor.Focus()
}

func (v *EditorView) Blur() {
	v.hostEditor.Blur()
}

func (v *EditorView) SetHost(host gohost.Host) {
	v.host = host
	v.hostEditor.SetValue(string(host.GetContent()))
}
