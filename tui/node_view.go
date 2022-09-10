package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gohost/log"
	"gohost/tui/form"
	"gohost/tui/keys"
)

var _ form.Form = (*NodeView)(nil)

type NodeView struct {
	model *Model
	*form.BaseForm
	nodeTypes *form.Choices
}

func NewNodeView(model *Model) *NodeView {
	// Text inputs
	nodeNameTextInput := form.NewTextInput()
	nodeNameTextInput.Prompt = "ID: "
	nodeNameTextInput.Focus(form.FocusFirstMode)

	descTextInput := form.NewTextInput()
	descTextInput.Prompt = "Description: "

	urlTextInput := form.NewTextInput()
	urlTextInput.Prompt = "Url: "

	// Node type choices
	nodeTypes := form.NewChoice([]list.DefaultItem{GroupNode, LocalHost, RemoteHost})

	nodeView := &NodeView{
		model:     model,
		BaseForm:  form.New(),
		nodeTypes: nodeTypes,
	}
	nodeView.WidgetStyle = lipgloss.NewStyle().PaddingBottom(1)
	nodeView.AddWidget(nodeNameTextInput)
	nodeView.AddWidget(descTextInput)
	nodeView.AddWidget(urlTextInput)
	nodeView.AddWidget(nodeTypes)

	return nodeView
}

func (v *NodeView) Init() tea.Cmd {
	v.model.setShortHelp(nodeViewState, keys.Arrows())
	return nil
}

func (v *NodeView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		v.SetSize(m.Width, m.Height)
		log.Debug(fmt.Sprintf("node view w %d h %d", m.Width, m.Height))
	case tea.KeyMsg:
		if v.model.state == nodeViewState {
			_, cmd = v.BaseForm.Update(msg)
		} else {
			return nil, nil
		}
	}
	cmds = append(cmds, cmd)
	return v, tea.Batch(cmds...)
}

func (v *NodeView) View() string {
	return v.BaseForm.View()
}
