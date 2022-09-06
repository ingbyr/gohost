package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gohost/gohost"
)

type EditorView struct {
	model      *Model
	hostEditor textarea.Model
	host       gohost.Host
	statusLine string
	saved      bool
	prevLen    int
}

func NewTextView(model *Model) *EditorView {
	hostEditor := textarea.New()
	hostEditor.ShowLineNumbers = true
	return &EditorView{
		model:      model,
		hostEditor: hostEditor,
		host:       nil,
		statusLine: "",
		saved:      true,
		prevLen:    0,
	}
}

func (v *EditorView) Init() tea.Cmd {
	return func() tea.Msg {
		// Display system host on start up
		v.SetHost(gohost.SysHost())
		return nil
	}
}

func (v *EditorView) Update(msg tea.Msg) []tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		v.hostEditor.SetHeight(m.Height - v.model.reservedHeight - 1)
		v.hostEditor.SetWidth(m.Width - v.model.groupView.groupList.Width())
	case tea.KeyMsg:
		if v.model.state == editorViewState {
			switch {
			case key.Matches(m, keys.Save):
				v.host.SetContent([]byte(v.hostEditor.Value()))
				err := gohost.GetService().UpdateHost(v.host)
				if err != nil {
					v.model.Log(err.Error())
				} else {
					v.SetSaved()
				}
			}
		} else {
			// Disable key
			msg = nil
		}
	}
	v.RefreshStatusLine()
	v.hostEditor, cmd = v.hostEditor.Update(msg)
	return append(cmds, cmd)
}

func (v *EditorView) View() string {
	return lipgloss.JoinVertical(lipgloss.Top, v.hostEditor.View(), v.statusLine)
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
	v.prevLen = v.hostEditor.Length()
}

func (v *EditorView) RefreshStatusLine() {
	v.statusLine = fmt.Sprintf("file: %s, saved: %t\n", v.host.GetName(), v.IsSaved())
}

func (v *EditorView) IsSaved() bool {
	saved := v.prevLen == v.hostEditor.Length()
	if !saved {
		v.saved = false
	}
	return v.saved && saved
}

func (v *EditorView) SetSaved() {
	v.prevLen = v.hostEditor.Length()
	v.saved = true
}
