package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"gohost/tui/styles"
	"strings"
)

type NodeView struct {
	model       *Model
	preFocusIdx int
	focusIdx    int
	inputs      []textinput.Model
}

func NewNodeView(model *Model) *NodeView {
	nodeNameTextInput := textinput.New()
	nodeNameTextInput.Prompt = "Name: "
	nodeNameTextInput.Focus()
	nodeNameTextInput.PromptStyle = styles.FocusedModel
	nodeNameTextInput.TextStyle = styles.FocusedModel

	descTextInput := textinput.New()
	descTextInput.Prompt = "Description: "

	view := &NodeView{
		model:       model,
		preFocusIdx: 0,
		focusIdx:    0,
		inputs:      []textinput.Model{nodeNameTextInput, descTextInput},
	}

	return view
}

func (v *NodeView) Init() tea.Cmd {
	return nil
}

func (v *NodeView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch m := msg.(type) {
	case tea.KeyMsg:
		if v.model.state == nodeViewState {
			switch {
			case key.Matches(m, keys.Enter, keys.Up, keys.Down):
				if key.Matches(m, keys.Enter, keys.Down) {
					cmds = append(cmds, v.setFocusInput(v.idxAfterFocusInput()))
				} else {
					cmds = append(cmds, v.setFocusInput(v.idxBeforeFocusInput()))
				}
			}
		} else {
			return nil, tea.Batch(cmds...)
		}
	}
	for i := range v.inputs {
		v.inputs[i], cmd = v.inputs[i].Update(msg)
		cmds = append(cmds, cmd)
	}
	return v, tea.Batch(cmds...)
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

func (v *NodeView) idxAfterFocusInput() int {
	id := v.focusIdx + 1
	if id >= len(v.inputs) {
		id = 0
	}
	return id
}

func (v *NodeView) idxBeforeFocusInput() int {
	id := v.focusIdx - 1
	if id < 0 {
		id = len(v.inputs) - 1
	}
	return id
}

func (v *NodeView) setFocusInput(idx int) tea.Cmd {
	v.preFocusIdx = v.focusIdx
	v.focusIdx = idx

	preInput := v.inputs[v.preFocusIdx]
	preInput.TextStyle = styles.None
	preInput.PromptStyle = styles.None
	preInput.Blur()
	v.inputs[v.preFocusIdx] = preInput

	input := v.inputs[v.focusIdx]
	input.TextStyle = styles.FocusedModel
	input.PromptStyle = styles.FocusedModel
	cmd := input.Focus()
	v.inputs[v.focusIdx] = input

	return cmd
}
