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
	hostEditor.CharLimit = 0
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
	km := v.hostEditor.KeyMap
	v.model.setShortHelp(editorViewState, []key.Binding{
		keys.Up,
		keys.Down,
		keys.Left,
		keys.Right,
		keys.Save,
		keys.Esc,
	})
	v.model.setFullHelp(editorViewState, [][]key.Binding{
		{keys.Up, keys.Down, keys.Left, keys.Right, keys.Save},
		{km.CharacterForward, km.CharacterBackward}, // TODO add all key map from textarea.KeyMap
	})
	return func() tea.Msg {
		// Display system host on start up
		v.SetHost(gohost.SysHost())
		return nil
	}
}

func (v *EditorView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		v.hostEditor.SetHeight(m.Height)
		v.hostEditor.SetWidth(m.Width)
	case tea.KeyMsg:
		if v.model.state == editorViewState {
			switch {
			case key.Matches(m, keys.Save):
				if v.host.IsEditable() {
					v.host.SetContent([]byte(v.hostEditor.Value()))
					err := gohost.GetService().UpdateHost(v.host)
					if err != nil {
						v.model.log(err.Error())
					} else {
						v.SetSaved()
					}
				} else {
					v.model.log("Can not edit this")
				}
			case key.Matches(m, keys.Esc):
				v.model.switchState(treeViewState)
			}
		} else {
			// Disable key
			msg = nil
		}
		v.statusLine = "hit key: " + m.String()
	}
	v.RefreshStatusLine()
	v.hostEditor, cmd = v.hostEditor.Update(msg)
	cmds = append(cmds, cmd)
	return v, tea.Batch(cmds...)
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
	v.hostEditor.Reset()
	v.hostEditor.SetValue(string(host.GetContent()))
	v.prevLen = v.hostEditor.Length()
}

func (v *EditorView) RefreshStatusLine() {
	v.statusLine = fmt.Sprintf("file: %s, saved: %t\n", v.host.Title(), v.IsSaved())
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
