package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"gohost/tui/styles"
	"strings"
)

type NodeView struct {
	model    *Model
	focusIdx int
	inputs   []textinput.Model
	labels   []string
}

func NewNodeView(model *Model) *NodeView {
	nodeNameTextInput := textinput.New()
	nodeNameTextInput.Prompt = "Name: "
	nodeNameTextInput.CursorStyle = styles.FocusedView.Copy()
	nodeNameTextInput.Focus()

	descTextInput := textinput.New()
	descTextInput.Prompt = "Description: "

	view := &NodeView{
		model:  model,
		inputs: []textinput.Model{nodeNameTextInput, descTextInput},
	}

	return view
}

func (v *NodeView) Init() tea.Cmd {
	return nil
}

func (v *NodeView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	for i := range v.inputs {
		v.inputs[i], cmd = v.inputs[i].Update(msg)
		cmds = append(cmds, cmd)
	}
	return v, cmd
}

func (v *NodeView) View() string {
	var b strings.Builder
	for i := range v.inputs {
		b.WriteString(v.inputs[i].View())
		if i < len(v.inputs)-1 {
			b.WriteString(cfg.LineBreak)
		}
	}
	return b.String()
}
