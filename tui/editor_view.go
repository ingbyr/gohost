package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gohost/gohost"
	"gohost/log"
	"gohost/tui/keys"
	"strings"
)

type EditorView struct {
	model         *Model
	hostEditor    textarea.Model
	hostNode      *gohost.TreeNode
	statusLine    string
	statusMsg     string
	saved         bool
	prevLen       int
	width, height int
}

func NewTextView(model *Model) *EditorView {
	hostEditor := textarea.New()
	hostEditor.ShowLineNumbers = true
	hostEditor.CharLimit = 0
	return &EditorView{
		model:      model,
		hostEditor: hostEditor,
		hostNode:   nil,
		statusLine: "",
		statusMsg:  "",
		saved:      true,
		prevLen:    0,
		width:      0,
		height:     0,
	}
}

func (v *EditorView) Init() tea.Cmd {
	km := v.hostEditor.KeyMap
	v.model.setShortHelp(StateEditorView, []key.Binding{
		keys.Up,
		keys.Down,
		keys.Left,
		keys.Right,
		keys.Save,
		keys.Esc,
	})
	v.model.setFullHelp(StateEditorView, [][]key.Binding{
		{keys.Up, keys.Down, keys.Left, keys.Right, keys.Save},
		{km.CharacterForward, km.CharacterBackward}, // TODO add all key map from textarea.KeyMap
	})
	v.SetHostNode(svc.SysHostNode)
	return nil
}

func (v *EditorView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		v.width, v.height = m.Width, m.Height
		v.hostEditor.SetHeight(m.Height - 2)
		v.hostEditor.SetWidth(m.Width)
		log.Debug(fmt.Sprintf("editor view w %d h %d", m.Width, m.Height))
	case AppliedNewHostMsg:
		v.SetHostNode(v.hostNode)
	case tea.KeyMsg:
		switch {
		case key.Matches(m, keys.Esc):
			return v, nil
		case key.Matches(m, keys.Save):
			host := v.Host()
			if host.IsEditable() {
				host.SetContent([]byte(v.hostEditor.Value()))
				svc.UpdateNode(v.hostNode)
				v.SetSaved()
			} else {
				v.statusMsg = "Can not edit this"
			}
		}
	}
	v.RefreshStatusLine()
	v.hostEditor, cmd = v.hostEditor.Update(msg)
	cmds = append(cmds, cmd)
	return v, tea.Batch(cmds...)
}

func (v *EditorView) View() string {
	var statusLine string
	if len(v.statusLine) > v.width {
		if v.width <= 3 {
			v.statusLine = strings.Repeat(" ", v.width)
		} else {
			statusLine = v.statusLine[:v.width-3] + "..."
		}
	} else {
		statusLine = v.statusLine
	}
	return lipgloss.JoinVertical(lipgloss.Right, v.hostEditor.View(), "", statusLine)
}

func (v *EditorView) Focus() {
	v.hostEditor.Focus()
}

func (v *EditorView) Blur() {
	v.hostEditor.Blur()
}

func (v *EditorView) Host() gohost.Host {
	if v.hostNode == nil {
		return nil
	}
	return v.hostNode.Node.(gohost.Host)
}

func (v *EditorView) SetHostNode(hostNode *gohost.TreeNode) {
	v.hostNode = hostNode
	v.hostEditor.Reset()
	v.hostEditor.SetValue(string(hostNode.Node.(gohost.Host).GetContent()))
	v.prevLen = v.hostEditor.Length()
}

func (v *EditorView) RefreshStatusLine() {
	if v.statusMsg == "" {
		v.statusLine = fmt.Sprintf("[file: %s] [saved: %t]", v.hostNode.Title(), v.IsSaved())
	} else {
		v.statusLine = fmt.Sprintf("[%s] [file: %s] [saved: %t]", v.statusMsg, v.hostNode.Title(), v.IsSaved())
	}
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
